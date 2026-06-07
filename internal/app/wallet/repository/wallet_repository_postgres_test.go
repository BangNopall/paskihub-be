package repository

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/BangNopall/paskihub-be/internal/infra/database"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func openWalletTestDatabase(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set; PostgreSQL locking test skipped")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	return db
}

func TestApproveTransactionSerializesConcurrentApprovalPostgres(t *testing.T) {
	db := openWalletTestDatabase(t)
	ctx := context.Background()
	if err := database.Migrate(db, nil); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	eventID := uuid.New()
	walletID := uuid.New()
	transactionID := uuid.New()
	t.Cleanup(func() {
		db.WithContext(ctx).Where("id = ?", transactionID).Delete(&entity.WalletTransaction{})
		db.WithContext(ctx).Where("id = ?", walletID).Delete(&entity.Wallet{})
		db.WithContext(ctx).Where("id = ?", eventID).Delete(&entity.Event{})
	})

	if err := db.WithContext(ctx).Create(&entity.Event{Id: eventID, UserId: uuid.New(), Name: "lock-test"}).Error; err != nil {
		t.Fatalf("create event: %v", err)
	}
	if err := db.WithContext(ctx).Create(&entity.Wallet{Id: walletID, EventId: eventID}).Error; err != nil {
		t.Fatalf("create wallet: %v", err)
	}
	if err := db.WithContext(ctx).Create(&entity.WalletTransaction{
		Id:       transactionID,
		WalletId: walletID,
		Type:     enums.TopUp,
		Amount:   1000,
		Status:   enums.Pending,
	}).Error; err != nil {
		t.Fatalf("create transaction: %v", err)
	}

	repo := NewWalletRepository(db)
	start := make(chan struct{})
	errs := make(chan error, 2)
	var wg sync.WaitGroup
	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			errs <- repo.ApproveTransaction(ctx, transactionID)
		}()
	}
	close(start)
	wg.Wait()
	close(errs)

	successes := 0
	badRequests := 0
	for err := range errs {
		switch err {
		case nil:
			successes++
		case domain.ErrBadRequest:
			badRequests++
		default:
			t.Fatalf("unexpected approval error: %v", err)
		}
	}
	if successes != 1 || badRequests != 1 {
		t.Fatalf("expected one success and one rejected replay, got success=%d rejected=%d", successes, badRequests)
	}

	var wallet entity.Wallet
	if err := db.WithContext(ctx).First(&wallet, "id = ?", walletID).Error; err != nil {
		t.Fatalf("reload wallet: %v", err)
	}
	if wallet.Saldo != 1000 {
		t.Fatalf("expected wallet saldo 1000, got %v", wallet.Saldo)
	}
}
