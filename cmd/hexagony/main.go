package main

import (
	"context"
	"hexagony/internal/app/api/rest"
	"hexagony/internal/app/repository/mysql"
	"log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mysqlRepository, err := mysql.NewMysqlRepository(ctx, "database_string")
	if err != nil {
		log.Println(err)
	}

	albumHandlers := rest.NewHandler(mysqlRepository)

	router := rest.Router(albumHandlers)

	rest.Server(ctx, router)
}
