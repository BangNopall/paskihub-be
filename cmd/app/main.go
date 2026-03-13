package main

import (
	"os"

	_ "github.com/BangNopall/paskihub-be/docs"
	"github.com/BangNopall/paskihub-be/internal/infra/database"
	"github.com/BangNopall/paskihub-be/internal/infra/env"
	"github.com/BangNopall/paskihub-be/internal/infra/server"
)

// @title Paskihub API
// @version 1.0
// @description This is the API documentation for Paskihub Backend.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3010
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	server := server.NewHttpServer()
	pgsqldb := database.NewPgsqlConn()

	database.Migrate(pgsqldb, os.Args)

	server.MountMiddlewares()
	server.MountRoutes(pgsqldb)
	server.RegistCustomValidation()
	server.Start(env.AppEnv.AppPort)
}