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
	// _ "github.com/go-sql-driver/mysql"

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

// @host  localhost:8000
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

	// enable or disable postgres based on the environment
	sslMode := "sslmode=disable"

	if envMode == "production" {
		sslMode = ""
	}

	// databaseURL := fmt.Sprintf(
	// 	"%s:%s@tcp(%s:%s)/%s?parseTime=true",
	// 	os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"),
	// 	os.Getenv("DB_PORT"), os.Getenv("DB_NAME"),
	// )

	// conn, err := sqlx.ConnectContext(ctx, "mysql", databaseURL) // mariadb uses the mysql driver
	// if err != nil {
	// 	clog.Fatal("mariadb failed to start")
	// }
	// defer conn.Close()

	// if err := conn.PingContext(ctx); err != nil {
	// 	clog.Fatal("could not ping mariadb database")
	// }

	databaseURL := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s %s",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), sslMode,
	)

	conn, err := sqlx.ConnectContext(ctx, "postgres", databaseURL) // mariadb uses the mysql driver
	if err != nil {
		clog.Info(err.Error())
		clog.Fatal("postgres failed to start")
	}
	defer conn.Close()

	if err := conn.PingContext(ctx); err != nil {
		clog.Fatal("could not ping postgres database")
	}

	router := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
		Debug:            os.Getenv("ENV_MODE") == "development",
	})

	router.Use(
		middleware.Timeout(time.Second*60),
		middleware.Recoverer,
		cmiddleware.LoggerMiddleware,
		render.SetContentType(render.ContentTypeJSON),
		cors.Handler,
	)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Welcome to Hexagony API")); err != nil {
			return
		}
	})

	router.Get("/docs/*", httpSwagger.WrapHandler)

	usersRepository := usersRepository.NewPostgresRepository(conn)
	usersController.NewUserHandler(router, usersRepository)

	albumsRepository := albumsRepository.NewPostgresRepository(conn)
	albumsController.NewAlbumHandler(router, albumsRepository)

	authRepository := authRepository.NewPostgresRepository(conn)
	authUseCase := authUseCase.NewAuthUsecase(authRepository)
	authController.NewAuthHandler(router, authUseCase)

	srv := &http.Server{
		Addr:              ":" + os.Getenv("PORT"),
		ReadTimeout:       time.Duration(time.Second * 5),
		ReadHeaderTimeout: time.Duration(time.Second * 5),
		WriteTimeout:      time.Duration(time.Second * 5),
		IdleTimeout:       time.Duration(time.Second * 20),
		Handler:           router,
	}

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
