package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/internal/infra/database/seeders"
	"github.com/BangNopall/paskihub-be/internal/infra/env"
	"github.com/BangNopall/paskihub-be/pkg/helpers/flag"
	"github.com/BangNopall/paskihub-be/pkg/log"
	"github.com/google/uuid"
)

const SEEDERS_FILE_PATH = "data/seeders/"
const SEEDERS_DEV_PATH = SEEDERS_FILE_PATH + "dev/"
const SEEDERS_PROD_PATH = SEEDERS_FILE_PATH + "prod/"

func getInterfaces() []interface{} {
	return []interface{}{
		&entity.User{},
		&entity.Institution{},
		&entity.Event{},
		&entity.EventLevel{},
		&entity.Team{},
		&entity.TeamMember{},
		&entity.Registration{},
		&entity.Wallet{},
		&entity.WalletTransaction{},
		&entity.TeamViolation{},
		&entity.ScoreSubCategory{},
		&entity.ScoreCategory{},
		&entity.Judge{},
		&entity.ViolationType{},
		&entity.GradeRule{},
		&entity.Score{},
		&entity.SystemSetting{},
		&entity.EventAward{},
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

func Migrate(db *gorm.DB, args []string) error {
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
		if err := db.Migrator().DropTable(getInterfaces()...); err != nil {
			return fmt.Errorf("drop tables: %w", err)
		}

	}

	log.Info(nil, "[PGSQL CONN][Migrate] Auto Migrating Tables")

	if err := db.Exec(`
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

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
	`).Error; err != nil {
		return fmt.Errorf("create database enums: %w", err)
	}

	if err := db.Exec(registrationStatusEnumSQL()).Error; err != nil {
		return fmt.Errorf("create registration status enum: %w", err)
	}

	if err := db.AutoMigrate(getInterfaces()...); err != nil {
		return fmt.Errorf("auto migrate tables: %w", err)
	}

	if err := EnsureDefaultSettings(db); err != nil {
		return err
	}
	return nil
}

func registrationStatusEnumSQL() string {
	return `
		DO $$ BEGIN
			CREATE TYPE registration_status AS ENUM (
				'WAITING',
				'DP_PAID',
				'FULL_PAID',
				'REJECTED',
				'KICKED'
			);
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;

		ALTER TYPE registration_status ADD VALUE IF NOT EXISTS 'KICKED';
	`
}

func EnsureDefaultSettings(db *gorm.DB) error {
	var count int64
	if err := db.Model(&entity.SystemSetting{}).Count(&count).Error; err != nil {
		return fmt.Errorf("count system settings: %w", err)
	}

	if count == 0 {
		log.Info(nil, "[PGSQL CONN][EnsureDefaultSettings] Seeding default system settings")
		defaultSettings := entity.SystemSetting{
			Id:            uuid.New(),
			CoinRate:      1000,
			ApprovalFee:   1,
			BankName:      "BCA",
			AccountNumber: "1234567890",
			AccountName:   "PaskiHub Indonesia",
		}
		if err := db.Create(&defaultSettings).Error; err != nil {
			return fmt.Errorf("create default system settings: %w", err)
		}
	}
	return nil
}

func Seeder(db *gorm.DB, flagVars *flag.Flag) {
	if !flagVars.Seeder {
		return
	}

	if flagVars.SeederModel == "User" {
		seeders.UserSeeder(db)
	} else if flagVars.SeederModel == "Demo" {
		seeders.DemoSeeder(db)
	} else if flagVars.SeederModel == "" {
		seeders.UserSeeder(db)
		// Add more seeders here
	} else {
		log.Error(nil, "[PGSQL CONN][Seeder] Seeder model not found")
	}
}
