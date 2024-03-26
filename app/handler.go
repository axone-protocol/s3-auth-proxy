package app

import (
	"encoding/json"

	"okp4/s3-auth-proxy/auth"

	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

func makeAuthenticateHandler(authenticator *auth.Authenticator) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		token, err := authenticator.Authenticate(ctx, ctx.Request.Body())
		if err != nil {
			ctx.Response.SetStatusCode(fasthttp.StatusForbidden)
			log.Info().Int("code", fasthttp.StatusForbidden).Err(err).Msg("🛑 VC authentication failed")
			return
		}

		body, err := json.Marshal(map[string]interface{}{
			"accessToken": token,
		})
		if err != nil {
			ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
			log.Info().Int("code", fasthttp.StatusInternalServerError).Err(err).Msg("🛑 Couldn't marshal response")
			return
		}

		ctx.Response.SetStatusCode(fasthttp.StatusOK)
		ctx.Response.SetBody(body)
		log.Info().Msg("✅ Authentication succeeded")
	}
}

func makeProxyHandler(s3Client *minio.Client, authenticator *auth.Authenticator) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		authHeader := string(ctx.Request.Header.Peek("Authorization"))

		bucket := ctx.UserValue("bucket").(string)
		filepath := ctx.UserValue("filepath").(string)
		logger := log.With().Str("bucket", bucket).Str("filepath", filepath).Logger()

		if len(authHeader) < 7 || authHeader[:6] != "Bearer" {
			ctx.Response.SetStatusCode(fasthttp.StatusUnauthorized)
			logger.Info().Int("code", fasthttp.StatusUnauthorized).Msg("🛑 Couldn't find bearer token")
			return
		}

		claims, err := authenticator.Authorize(authHeader[7:])
		if err != nil {
			ctx.Response.SetStatusCode(fasthttp.StatusUnauthorized)
			ctx.Response.SetBody([]byte(err.Error()))
			logger.Info().Int("code", fasthttp.StatusUnauthorized).Err(err).Msg("🛑 Invalid jwt")
			return
		}

		logger = logger.With().Str("aud", claims.Audience).Str("jti", claims.Id).Logger()
		obj, err := s3Client.GetObject(ctx, bucket, filepath, minio.GetObjectOptions{})
		if err != nil {
			ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
			logger.Info().Int("code", fasthttp.StatusInternalServerError).Err(err).Msg("😿 Could not proxy the request")
			return
		}

		info, err := obj.Stat()
		if err != nil {
			ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
			logger.Info().Int("code", fasthttp.StatusInternalServerError).Err(err).Msg("😿 Could not proxy the request")
			return
		}

		ctx.Response.SetStatusCode(fasthttp.StatusOK)
		ctx.Response.Header.SetContentType(info.ContentType)
		ctx.Response.SetBodyStream(obj, -1)

		logger.Info().Int("code", fasthttp.StatusOK).Msg("✅ Proxying request")
	}
}
