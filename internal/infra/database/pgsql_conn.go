package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/internal/infra/env"
	"github.com/BangNopall/paskihub-be/pkg/helpers/flag"
	"github.com/BangNopall/paskihub-be/pkg/log"
)

const SEEDERS_FILE_PATH = "data/seeders/"
const SEEDERS_DEV_PATH = SEEDERS_FILE_PATH + "dev/"
const SEEDERS_PROD_PATH = SEEDERS_FILE_PATH + "prod/"

func getInterfaces() []interface{} {
	return []interface{}{
		&entity.User{},
	}
}

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

func Migrate(db *gorm.DB, args []string) {
	if flag.FlagVars.Fresh {

		if env.AppEnv.AppEnv == "production" {
			var choice string
			fmt.Print("Application is on production. Are you sure you want to do fresh migration ? (y/n): ")
			fmt.Scan(&choice)

			if choice != "y" {
				fmt.Print("Exiting...\n")
				os.Exit(0)
			}
		}

		log.Info(nil, "[PGSQL CONN][Migrate] Dropping All Tables")
		db.Migrator().DropTable(getInterfaces()...)

	}

	log.Info(nil, "[PGSQL CONN][Migrate] Auto Migrating Tables")

	db.AutoMigrate(getInterfaces()...)
}
