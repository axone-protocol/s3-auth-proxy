package app

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/okp4/s3-auth-proxy/auth"

	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
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
	server := configure(a.s3Client, a.authenticator)
	ln, err := net.Listen("tcp4", a.listenAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("❌ Couldn't start server")
	}

	listenErr := make(chan error, 1)
	go func() {
		log.Info().Str("listenAddr", a.listenAddr).Msg("🔥 Listening")
		listenErr <- server.Serve(ln)
	}()

	kill := make(chan os.Signal, 1)
	signal.Notify(kill, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-listenErr:
			if err != nil {
				log.Fatal().Err(err).Msg("❌ Listening error")
			}
			return
		case <-kill:
			log.Info().Msg("🧯 Shutting down")
			if err := ln.Close(); err != nil {
				log.Fatal().Err(err).Msg("❌ Couldn't stop listener")
			}
		}
	}
}
