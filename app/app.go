package app

import (
	"github.com/rs/zerolog/log"
	"okp4/minio-auth-plugin/auth"
)

type AuthApp struct {
	listenAddr    string
	authenticator *auth.Authenticator
}

func New(listenAddr string, authenticator *auth.Authenticator) *AuthApp {
	return &AuthApp{
		listenAddr:    listenAddr,
		authenticator: authenticator,
	}
}

func (a *AuthApp) Start() {
	log.Info().Str("listenAddr", a.listenAddr).Msg("🔥 Listening")
	log.Fatal().Err(configure(a.authenticator).ListenAndServe(a.listenAddr)).Msg("🛑 Shutting down")
}
