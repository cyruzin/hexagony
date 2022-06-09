package main

import (
	"context"
	"fmt"

	"hexagony/internal/app/config"
	albumController "hexagony/internal/app/modules/albums/infra/controller"
	albumRepository "hexagony/internal/app/modules/albums/repository/mysql"
	sharedMiddleware "hexagony/internal/app/modules/shared/infra/middleware"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "hexagony/docs"

	_ "github.com/go-sql-driver/mysql"
)

// @title           Hexagony API
// @version         1.0
// @description     Clean architecture example in Golang.

// @contact.name   Cyro Dubeux
// @contact.url    https://github.com/cyruzin/hexagony/issues/new
// @contact.email  xorycx@gmail.com

// @license.name  MIT
// @license.url   https://github.com/cyruzin/hexagony/blob/master/LICENSE

// @host      localhost:8000
// @BasePath  /
func main() {
	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	if cfg.EnvMode == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Debug().Msg("running in development mode")
	} else {
		log.Info().Msg("running in production mode")
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	}

	databaseURL := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPass, cfg.DBHost,
		cfg.DBPort, cfg.DBName,
	)

	conn, err := sqlx.ConnectContext(ctx, "mysql", databaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("mysql failed to start")
	}
	defer conn.Close()

	if err := conn.PingContext(ctx); err != nil {
		log.Fatal().
			Err(err).
			Stack().
			Str("database", conn.DriverName()).
			Msg("could not ping the database")
	}

	router := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	router.Use(
		cors.Handler,
		middleware.Timeout(cfg.MiddlewareTimeOut),
		render.SetContentType(render.ContentTypeJSON),
		sharedMiddleware.LoggerMiddleware,
	)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hexagony v1.0"))
	})

	router.Get("/docs/*", httpSwagger.WrapHandler)

	albumRepository := albumRepository.NewMysqlRepository(conn)
	albumController.NewAlbumHandler(router, albumRepository)

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

		log.Info().Msg("shutting down the server...")
		if err := srv.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("server failed to shutdown")
		}
		close(idleConnsClosed)
	}()

	log.Info().Msgf("listening on port: %s", cfg.Port)
	log.Info().Msg("you're good to go! :)")

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Error().Err(err).Msg("server failed to start")
	}

	<-idleConnsClosed
}
