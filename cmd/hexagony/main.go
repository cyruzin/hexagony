package main

import (
	"context"
	"fmt"
	"hexagony/internal/app/api/rest"
	"hexagony/internal/app/config"
	"hexagony/internal/app/repository/mysql"
	"log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	databaseURL := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPass, cfg.DBHost,
		cfg.DBPort, cfg.DBName,
	)

	mysqlRepository, err := mysql.NewMysqlRepository(ctx, databaseURL)
	if err != nil {
		log.Fatal(err.Error())
	}

	albumHandlers := rest.NewHandler(mysqlRepository)
	router := rest.Router(albumHandlers)

	rest.Server(ctx, router)
}
