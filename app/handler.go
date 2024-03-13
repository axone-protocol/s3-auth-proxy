package app

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"okp4/minio-auth-plugin/auth"
)

func makeAuthenticateHandler(authenticator *auth.Authenticator) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		log.Info().Msg("Authentication requested")
		ctx.Response.SetStatusCode(fasthttp.StatusUnauthorized)
	}
}

func makeProxyHandler(minioClient *minio.Client, authenticator *auth.Authenticator) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		authHeader := string(ctx.Request.Header.Peek("Authorization"))

		bucket := ctx.UserValue("bucket").(string)
		filepath := ctx.UserValue("filepath").(string)
		logger := log.With().Str("bucket", bucket).Str("filepath", filepath).Logger()

		if len(authHeader) < 7 || authHeader[:6] != "Bearer" {
			ctx.Response.SetStatusCode(fasthttp.StatusUnauthorized)
			logger.Info().Int("code", fasthttp.StatusUnauthorized).Msg("Couldn't find bearer token")
			return
		}

		if err := authenticator.Authorize(authHeader[7:]); err != nil {
			ctx.Response.SetStatusCode(fasthttp.StatusUnauthorized)
			logger.Info().Int("code", fasthttp.StatusUnauthorized).Msg("Invalid jwt")
			return
		}

		obj, err := minioClient.GetObject(context.Background(), bucket, filepath, minio.GetObjectOptions{})
		if err != nil {
			ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
			logger.Info().Int("code", fasthttp.StatusInternalServerError).Err(err).Msg("Could not proxy the request")
			return
		}

		info, err := obj.Stat()
		if err != nil {
			ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
			logger.Info().Int("code", fasthttp.StatusInternalServerError).Err(err).Msg("Could not proxy the request")
			return
		}

		ctx.Response.SetStatusCode(fasthttp.StatusOK)
		ctx.Response.Header.SetContentType(info.ContentType)
		ctx.Response.SetBodyStream(obj, -1)

		logger.Info().Int("code", fasthttp.StatusOK).Msg("Proxying request")

	}
}
