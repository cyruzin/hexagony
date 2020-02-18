package rest

import (
	"context"
	"hexagony/internal/app/config"
	"net/http"
	"os"
	"os/signal"

	"github.com/rs/zerolog/log"
)

// Server initiates the server.
func Server(
	ctx context.Context,
	cfg *config.Config,
	router http.Handler,
) {
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		ReadTimeout:       cfg.ReadTimeOut,
		ReadHeaderTimeout: cfg.ReadHeaderTimeOut,
		WriteTimeout:      cfg.WriteTimeOut,
		IdleTimeout:       cfg.IdleTimeOut,
		Handler:           router,
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		gracefulStop := make(chan os.Signal, 1)
		signal.Notify(gracefulStop, os.Interrupt)
		<-gracefulStop

		log.Info().Msg("Shutting down the server...")
		if err := srv.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("Server failed to shutdown")
		}
		close(idleConnsClosed)
	}()

	log.Info().Msgf("Listening on port: %s", cfg.Port)
	log.Info().Msg("You're good to go! :)")

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Error().Err(err).Msg("Server failed to start")
	}

	<-idleConnsClosed
}
