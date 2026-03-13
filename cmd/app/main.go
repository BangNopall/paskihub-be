package main

import (
	"os"

	_ "github.com/BangNopall/paskihub-be/docs"
	"github.com/BangNopall/paskihub-be/internal/infra/database"
	"github.com/BangNopall/paskihub-be/internal/infra/env"
	"github.com/BangNopall/paskihub-be/internal/infra/server"
)

// @title						Paskihub API
// @version					1.0
// @description				This is Paskihub API Documentation
// @host						localhost:3010
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
	server := server.NewHttpServer()
	pgsqldb := database.NewPgsqlConn()

	database.Migrate(pgsqldb, os.Args)

	server.MountMiddlewares()
	server.MountRoutes(pgsqldb)
	server.RegistCustomValidation()
	server.Start(env.AppEnv.AppPort)
}
