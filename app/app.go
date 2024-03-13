package app

import (
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
	"okp4/minio-auth-plugin/auth"
)

type AuthApp struct {
	listenAddr    string
	minioClient   *minio.Client
	authenticator *auth.Authenticator
}

func New(listenAddr string, minioClient *minio.Client, authenticator *auth.Authenticator) *AuthApp {
	return &AuthApp{
		listenAddr:    listenAddr,
		minioClient:   minioClient,
		authenticator: authenticator,
	}
}

func (a *AuthApp) Start() {
	log.Info().Str("listenAddr", a.listenAddr).Msg("🔥 Listening")
	log.Fatal().Err(configure(a.minioClient, a.authenticator).ListenAndServe(a.listenAddr)).Msg("🛑 Shutting down")
}
