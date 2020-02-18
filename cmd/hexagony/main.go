package main

import (
	"context"
	"fmt"
	"hexagony/internal/app/api/rest"
	"hexagony/internal/app/config"
	"hexagony/internal/app/repository/mysql"

	"github.com/rs/zerolog"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	databaseURL := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPass, cfg.DBHost,
		cfg.DBPort, cfg.DBName,
	)

	mysqlRepository := mysql.NewMysqlRepository(databaseURL)

	albumHandlers := rest.NewHandler(mysqlRepository)
	router := rest.Router(albumHandlers)

	rest.Server(ctx, cfg, router)
}
