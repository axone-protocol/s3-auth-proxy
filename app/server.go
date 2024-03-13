package app

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"okp4/minio-auth-plugin/auth"
)

func configure(authenticator *auth.Authenticator) *fasthttp.Server {
	r := router.New()
	r.POST("/auth", makeAuthenticateHandler(authenticator))
	r.POST("/authz", makeAuthPluginHandler(authenticator))

	return &fasthttp.Server{
		Handler: r.Handler,
	}
}
