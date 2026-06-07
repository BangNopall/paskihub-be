package database

import (
	"strings"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func disconnectedDatabase(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(
		postgres.New(postgres.Config{
			DSN:                  "host=127.0.0.1 port=1 user=invalid dbname=invalid sslmode=disable",
			PreferSimpleProtocol: true,
		}),
		&gorm.Config{DisableAutomaticPing: true},
	)
	if err != nil {
		t.Fatalf("construct disconnected database: %v", err)
	}
	return db
}

func TestRegistrationStatusMigrationIncludesKicked(t *testing.T) {
	sql := registrationStatusEnumSQL()
	if !strings.Contains(sql, "'KICKED'") {
		t.Fatalf("registration enum migration does not include KICKED: %s", sql)
	}
	if !strings.Contains(sql, "ADD VALUE IF NOT EXISTS") {
		t.Fatalf("registration enum migration is not safe for existing databases: %s", sql)
	}
}

func TestMigrateReturnsDatabaseErrors(t *testing.T) {
	if err := Migrate(disconnectedDatabase(t), nil); err == nil {
		t.Fatal("expected migration error from unavailable database")
	}
}

func TestEnsureDefaultSettingsReturnsDatabaseErrors(t *testing.T) {
	if err := EnsureDefaultSettings(disconnectedDatabase(t)); err == nil {
		t.Fatal("expected default setting error from unavailable database")
	}
}
