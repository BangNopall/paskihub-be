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
		&entity.Event{},
		&entity.EventLevel{},
		&entity.Registration{},
		&entity.Wallet{},
		&entity.WalletTransaction{},
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

	db.Exec(`
		DO $$ BEGIN
			CREATE TYPE role AS ENUM (
				'ADMIN',
				'ORGANIZER',
				'PESERTA'
			);
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;

		DO $$ BEGIN
			CREATE TYPE event_status AS ENUM (
				'DRAFT',
				'OPEN',
				'CLOSED',
				'ARCHIVED'
			);
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;

		DO $$ BEGIN
			CREATE TYPE institution_type AS ENUM (
				'SD',
				'SMP',
				'SMA',
				'PURNA',
				'UMUM'
			);
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;

		DO $$ BEGIN
			CREATE TYPE team_type AS ENUM (
				'PASUKAN',
				'DANPAS',
				'OFFICIAL',
				'PELATIH'
			);
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;

		DO $$ BEGIN
			CREATE TYPE registration_status AS ENUM (
				'WAITING',
				'DP_PAID',
				'FULL_PAID',
				'REJECTED'
			);
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;

		DO $$ BEGIN
			CREATE TYPE wallet_type AS ENUM (
				'TOPUP',
				'WITHDRAW'
			);
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;

		DO $$ BEGIN
			CREATE TYPE transaction_status AS ENUM (
				'PENDING',
				'APPROVE',
				'REJECTED'
			);
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
	`)

	db.AutoMigrate(getInterfaces()...)
}
