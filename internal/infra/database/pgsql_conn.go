package database

import (
	"fmt"
	// "os"
	// "reflect"
	// "strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/BangNopall/paskihub-be/internal/infra/env"
	"github.com/BangNopall/paskihub-be/pkg/log"
)

const SEEDERS_FILE_PATH = "data/seeders/"
const SEEDERS_DEV_PATH = SEEDERS_FILE_PATH + "dev/"
const SEEDERS_PROD_PATH = SEEDERS_FILE_PATH + "prod/"

func NewPgsqlConn() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		env.AppEnv.DBHost,
		env.AppEnv.DBUser,
		env.AppEnv.DBPass,
		env.AppEnv.DBName,
		env.AppEnv.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})

	if err != nil {
		log.Fatal(log.LogInfo{
			"error": err.Error(),
		}, "[PGSQL CONN][NewPgsqlConn] Failed to connect to database")
	}

	return db
}