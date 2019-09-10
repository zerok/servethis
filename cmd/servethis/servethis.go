package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
	"github.com/zerok/servethis/pkg/server"
)

func main() {
	ctx := context.Background()
	var addr string
	pflag.StringVar(&addr, "addr", "127.0.0.1:9980", "Address to listen on")
	pflag.Parse()

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.InfoLevel).With().Timestamp().Logger()
	ctx = logger.WithContext(ctx)

	wd, err := os.Getwd()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to determine current working directory.")
	}

	logger.Info().Msgf("Serving content from %s", wd)

	handler := server.New(ctx, wd)
	server := http.Server{
		ReadTimeout: time.Second * 2,
		Addr:        addr,
		Handler:     handler,
	}

	logger.Info().Msgf("Starting listener on %s", addr)
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to start listener.")
	}
}
