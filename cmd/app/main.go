package main

import (
	"os"

	_ "github.com/BangNopall/paskihub-be/docs"
	"github.com/BangNopall/paskihub-be/internal/infra/database"
	"github.com/BangNopall/paskihub-be/internal/infra/env"
	"github.com/BangNopall/paskihub-be/internal/infra/server"
	"github.com/BangNopall/paskihub-be/pkg/helpers/flag"
	"github.com/BangNopall/paskihub-be/pkg/log"
)

// @title						Paskihub API
// @version					1.0
// @description				This is Paskihub API Documentation
// @Security ApiKeyAuth
// @host						localhost:3080
// @schemes					http
// @BasePath 				/
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				JWT Bearer token. Format: Bearer {token}
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						x-api-key
// @description				API Key for all endpoints. Format: Key {api_key}. Example: Key abc123
func main() {
	if err := flag.Parse(os.Args[1:]); err != nil {
		log.Fatal(log.LogInfo{
			"error": err.Error(),
		}, "[APP][main] invalid command flags")
	}

	server := server.NewHttpServer()
	pgsqldb := database.NewPgsqlConn()

	if err := database.Migrate(pgsqldb, os.Args); err != nil {
		log.Fatal(log.LogInfo{
			"error": err.Error(),
		}, "[APP][main] database migration failed")
	}
	database.Seeder(pgsqldb, flag.FlagVars)

	server.MountMiddlewares()
	if err := server.MountRoutes(pgsqldb); err != nil {
		log.Fatal(log.LogInfo{
			"error": err.Error(),
		}, "[APP][main] failed to mount routes")
	}
	server.RegistCustomValidation()
	server.Start(env.AppEnv.AppPort)
}
