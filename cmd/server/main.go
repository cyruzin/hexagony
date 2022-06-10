package main

import (
	"context"
	"fmt"
	"time"

	albumController "hexagony/internal/app/albums/infra/controller"
	albumRepository "hexagony/internal/app/albums/repository/mysql"
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

// @title        Hexagony API
// @version      1.0
// @description  Clean architecture example in Golang.

// @contact.name   Cyro Dubeux
// @contact.url    https://github.com/cyruzin/hexagony/issues/new
// @contact.email  xorycx@gmail.com

// @license.name  MIT
// @license.url   https://github.com/cyruzin/hexagony/blob/master/LICENSE

// @host      localhost:8000
// @BasePath  /
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	if os.Getenv("ENV_MODE") == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Debug().Msg("running in development mode")
	} else {
		log.Info().Msg("running in production mode")
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	}

	databaseURL := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"), os.Getenv("DB_NAME"),
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
		render.SetContentType(render.ContentTypeJSON),
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Timeout(time.Second*60),
	)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hexagony v1.0")); err != nil {
			return
		}
	})

	router.Get("/docs/*", httpSwagger.WrapHandler)

	albumRepository := albumRepository.NewMysqlRepository(conn)
	albumController.NewAlbumHandler(router, albumRepository)

	srv := &http.Server{
		Addr:              ":" + os.Getenv("PORT"),
		ReadTimeout:       time.Duration(time.Second * 5),
		ReadHeaderTimeout: time.Duration(time.Second * 5),
		WriteTimeout:      time.Duration(time.Second * 5),
		IdleTimeout:       time.Duration(time.Second * 60),
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

	log.Info().Msgf("listening on port: %s", os.Getenv("PORT"))
	log.Info().Msg("you're good to go! :)")

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Error().Err(err).Msg("server failed to start")
	}

	<-idleConnsClosed
}
