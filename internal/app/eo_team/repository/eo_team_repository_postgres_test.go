package repository

import (
	"context"
	"testing"

	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(
		&entity.User{},
		&entity.Event{},
		&entity.EventLevel{},
		&entity.Team{},
		&entity.Registration{},
		&entity.Wallet{},
		&entity.WalletTransaction{},
	)
	return db
}

func TestApproveRegistrationDebitsWalletOncePostgres(t *testing.T) {
	db := setupTestDB()
	repo := NewEOTeamRepository(db)
	ctx := context.Background()

	eventID := uuid.New()
	walletID := uuid.New()
	registrationID := uuid.New()
	levelID := uuid.New()

	db.Create(&entity.Event{Id: eventID, Status: "OPEN"})
	db.Create(&entity.Wallet{Id: walletID, EventId: eventID, Saldo: 0, CoinBalance: 100})
	db.Create(&entity.EventLevel{Id: levelID, EventId: eventID})
	db.Create(&entity.Registration{
		Id:             registrationID,
		EventLevelId:   levelID,
		PaymentStatus:  enums.Waiting,
	})

	err := repo.ApproveRegistration(ctx, eventID, registrationID, 1000, 10, 1000, enums.FullPaid)
	assert.NoError(t, err)

	var wallet entity.Wallet
	db.First(&wallet, "id = ?", walletID)
	assert.Equal(t, float64(90), wallet.CoinBalance)
	assert.Equal(t, float64(-10000), wallet.Saldo)

	var tx entity.WalletTransaction
	db.First(&tx, "wallet_id = ?", walletID)
	assert.Equal(t, float64(-10000), tx.Amount)
	assert.Equal(t, enums.Withdraw, tx.Type)
}
