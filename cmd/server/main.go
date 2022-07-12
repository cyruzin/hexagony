package main

import (
	"context"
	"fmt"
	"time"

	albumsController "hexagony/internal/albums/controller"
	albumsRepository "hexagony/internal/albums/repository/postgres"
	usersController "hexagony/internal/users/controller"
	usersRepository "hexagony/internal/users/repository/postgres"
	"hexagony/pkg/clog"

	authController "hexagony/internal/auth/controller"
	authRepository "hexagony/internal/auth/repository/postgres"
	authUseCase "hexagony/internal/auth/usecase"

	"net/http"
	"os"
	"os/signal"

	cmiddleware "hexagony/internal/shared/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "hexagony/docs"

	_ "github.com/lib/pq"
)

// @title        Hexagony API
// @version      1.0
// @description  Clean architecture example in Golang.

// @contact.name   Cyro Dubeux
// @contact.url    https://github.com/cyruzin/hexagony/issues/new
// @contact.email  xorycx@gmail.com

// @license.name  MIT
// @license.url   https://github.com/cyruzin/hexagony/blob/master/LICENSE

// @host  hexagony.herokuapp.com
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	envMode := os.Getenv("ENV_MODE")

	if envMode == "development" {
		clog.UseConsoleOutput()
		clog.Debug("running in development mode")
	} else {
		clog.Info("running in production mode")
	}

	// postgres url string
	databaseURL := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"),
	)

	if envMode == "production" { // get url from heroku if production
		databaseURL = os.Getenv("DATABASE_URL")
	}

	// connecting to postgres
	conn, err := sqlx.ConnectContext(ctx, "postgres", databaseURL)
	if err != nil {
		clog.Info(err.Error())
		clog.Fatal("postgres failed to start")
	}
	defer conn.Close()

	if err := conn.PingContext(ctx); err != nil {
		clog.Fatal("could not ping postgres database")
	}

	router := chi.NewRouter()

	// enabling CORS
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
		Debug:            envMode == "development",
	})

	// middlewares
	router.Use(
		middleware.Timeout(time.Second*60),
		middleware.Recoverer,
		cmiddleware.LoggerMiddleware,
		render.SetContentType(render.ContentTypeJSON),
		cors.Handler,
	)

	// root page
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Welcome to Hexagony API")); err != nil {
			return
		}
	})

	// swagger documentation end-point
	router.Get("/docs/*", httpSwagger.WrapHandler)

	// domain instances
	usersRepository := usersRepository.NewPostgresRepository(conn)
	usersController.NewUserHandler(router, usersRepository)

	albumsRepository := albumsRepository.NewPostgresRepository(conn)
	albumsController.NewAlbumHandler(router, albumsRepository)

	authRepository := authRepository.NewPostgresRepository(conn)
	authUseCase := authUseCase.NewAuthUsecase(authRepository)
	authController.NewAuthHandler(router, authUseCase)

	// server configuration and timeouts
	srv := &http.Server{
		Addr:              ":" + os.Getenv("PORT"),
		ReadTimeout:       time.Duration(time.Second * 5),
		ReadHeaderTimeout: time.Duration(time.Second * 5),
		WriteTimeout:      time.Duration(time.Second * 5),
		IdleTimeout:       time.Duration(time.Second * 20),
		Handler:           router,
	}

	// graceful shutdown
	idleConnsClosed := make(chan struct{})

	go func() {
		gracefulStop := make(chan os.Signal, 1)
		signal.Notify(gracefulStop, os.Interrupt)
		<-gracefulStop

		clog.Info("shutting down the server...")
		if err := srv.Shutdown(ctx); err != nil {
			clog.Error(err, "server failed to shutdown")
		}
		close(idleConnsClosed)
	}()

	clog.Info("listening on port: " + os.Getenv("PORT"))
	clog.Info("you're good to go! :)")

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		clog.Error(err, "server failed to start")
	}

	<-idleConnsClosed
}
