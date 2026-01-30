package main

import (
	"os"

	"github.com/BangNopall/paskihub-be/internal/infra/database"
	"github.com/BangNopall/paskihub-be/internal/infra/env"
	"github.com/BangNopall/paskihub-be/internal/infra/server"
)

func main() {
	server := server.NewHttpServer()
	pgsqldb := database.NewPgsqlConn()

	database.Migrate(pgsqldb, os.Args)

	server.MountMiddlewares()
	server.MountRoutes(pgsqldb)
	server.RegistCustomValidation()
	server.Start(env.AppEnv.AppPort)
}