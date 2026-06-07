package database

import (
	"context"
	"os"
	"testing"

	"github.com/BangNopall/paskihub-be/domain/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestMigrationCreatesKickedEnumAndDefaultSettingPostgres(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set; PostgreSQL migration test skipped")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := Migrate(db, nil); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	var count int64
	err = db.WithContext(context.Background()).Raw(`
		SELECT COUNT(*)
		FROM pg_enum pe
		JOIN pg_type pt ON pt.oid = pe.enumtypid
		WHERE pt.typname = 'registration_status' AND pe.enumlabel = 'KICKED'
	`).Scan(&count).Error
	if err != nil {
		t.Fatalf("query registration enum: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected KICKED enum value, got count %d", count)
	}

	count = 0
	if err := db.Model(&entity.SystemSetting{}).Count(&count).Error; err != nil {
		t.Fatalf("count default settings: %v", err)
	}
	if count == 0 {
		t.Fatal("default system setting was not created")
	}
}
