package main

import (
	"context"
	"hexagony/internal/app/api/rest"
	"hexagony/internal/app/repository/mysql"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mysqlRepository, err := mysql.NewMysqlRepository(ctx, "database_string")
	if err != nil {
		log.Println(err)
	}

	albumHandlers := rest.NewHandler(mysqlRepository)

	router := chi.NewRouter()
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(middleware.Logger)

	router.Get("/album", albumHandlers.Index)
	router.Get("/album/{id}", albumHandlers.Show)
	router.Post("/album", albumHandlers.Store)
	router.Put("/album/{id}", albumHandlers.Update)
	router.Delete("/album/{id}", albumHandlers.Delete)

	http.ListenAndServe(":8000", router)
}
