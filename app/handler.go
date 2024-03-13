package app

import (
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

func makeAuthPluginHandler(authenticator *auth.Authenticator) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		log.Info().Msg("Authorization requested")
		ctx.Response.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.Response.SetBody([]byte("{\"result\":true}"))
	}
}
