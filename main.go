package main

import (
	"os"

	"github.com/axone-protocol/s3-auth-proxy/cmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cmd.Execute()
}
