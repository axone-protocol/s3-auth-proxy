package app

import (
	"github.com/okp4/s3-auth-proxy/auth"

	"github.com/fasthttp/router"
	"github.com/minio/minio-go/v7"
	"github.com/valyala/fasthttp"
)

func configure(minioClient *minio.Client, authenticator *auth.Authenticator) *fasthttp.Server {
	r := router.New()
	r.POST("/auth", makeAuthenticateHandler(authenticator))
	r.GET("/{bucket}/{filepath:*}", makeProxyHandler(minioClient, authenticator))

	return &fasthttp.Server{
		Handler: r.Handler,
	}
}
