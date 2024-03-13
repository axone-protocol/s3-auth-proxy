package app

import (
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
	"okp4/s3-auth-proxy/auth"
)

type AuthApp struct {
	listenAddr    string
	s3Client      *minio.Client
	authenticator *auth.Authenticator
}

func New(listenAddr string, s3Client *minio.Client, authenticator *auth.Authenticator) *AuthApp {
	return &AuthApp{
		listenAddr:    listenAddr,
		s3Client:      s3Client,
		authenticator: authenticator,
	}
}

func (a *AuthApp) Start() {
	log.Info().Str("listenAddr", a.listenAddr).Msg("🔥 Listening")
	log.Fatal().Err(configure(a.s3Client, a.authenticator).ListenAndServe(a.listenAddr)).Msg("🛑 Shutting down")
}
