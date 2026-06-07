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

func TestApproveRegistrationDebitsWalletOncePostgres(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set; PostgreSQL team approval test skipped")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := database.Migrate(db, nil); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	ctx := context.Background()
	organizerID := uuid.New()
	participantID := uuid.New()
	institutionID := uuid.New()
	eventID := uuid.New()
	eventLevelID := uuid.New()
	teamID := uuid.New()
	registrationID := uuid.New()
	walletID := uuid.New()

	t.Cleanup(func() {
		db.WithContext(ctx).Where("wallet_id = ?", walletID).Delete(&entity.WalletTransaction{})
		db.WithContext(ctx).Where("id = ?", registrationID).Delete(&entity.Registration{})
		db.WithContext(ctx).Where("id = ?", walletID).Delete(&entity.Wallet{})
		db.WithContext(ctx).Where("id = ?", teamID).Delete(&entity.Team{})
		db.WithContext(ctx).Where("id = ?", eventLevelID).Delete(&entity.EventLevel{})
		db.WithContext(ctx).Where("id = ?", eventID).Delete(&entity.Event{})
		db.WithContext(ctx).Where("id = ?", institutionID).Delete(&entity.Institution{})
		db.WithContext(ctx).Where("id IN ?", []uuid.UUID{organizerID, participantID}).Delete(&entity.User{})
	})

	fixtures := []interface{}{
		&entity.User{Id: organizerID, Email: organizerID.String() + "@test.local", Password: "x", Role: enums.Organizer},
		&entity.User{Id: participantID, Email: participantID.String() + "@test.local", Password: "x", Role: enums.Peserta},
		&entity.Institution{Id: institutionID, UserId: participantID, Name: "Test Institution", InstitutionType: enums.SMA},
		&entity.Event{Id: eventID, UserId: organizerID, Name: "Concurrent Approval"},
		&entity.EventLevel{Id: eventLevelID, EventId: eventID, Name: "SMA"},
		&entity.Team{Id: teamID, InstiId: institutionID, Name: "Test Team"},
		&entity.Registration{Id: registrationID, TeamId: teamID, EventLevelId: eventLevelID, PaymentStatus: enums.Waiting},
		&entity.Wallet{Id: walletID, EventId: eventID, Saldo: 5000},
	}
	for _, fixture := range fixtures {
		if err := db.WithContext(ctx).Create(fixture).Error; err != nil {
			t.Fatalf("create fixture %T: %v", fixture, err)
		}
	}

	repo := NewEOTeamRepository(db)
	start := make(chan struct{})
	errs := make(chan error, 2)
	var wg sync.WaitGroup
	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			errs <- repo.ApproveRegistration(ctx, eventID, registrationID, 1000, enums.FullPaid)
		}()
	}
	close(start)
	wg.Wait()
	close(errs)

	successes := 0
	replays := 0
	for err := range errs {
		switch err {
		case nil:
			successes++
		case domain.ErrBadRequest:
			replays++
		default:
			t.Fatalf("unexpected approval error: %v", err)
		}
	}
	if successes != 1 || replays != 1 {
		t.Fatalf("expected one success and one replay rejection, got success=%d replay=%d", successes, replays)
	}

	var wallet entity.Wallet
	if err := db.WithContext(ctx).First(&wallet, "id = ?", walletID).Error; err != nil {
		t.Fatalf("reload wallet: %v", err)
	}
	if wallet.Saldo != 4000 {
		t.Fatalf("expected saldo 4000 after one debit, got %v", wallet.Saldo)
	}

	var withdrawals int64
	if err := db.WithContext(ctx).
		Model(&entity.WalletTransaction{}).
		Where("wallet_id = ? AND type = ?", walletID, enums.Withdraw).
		Count(&withdrawals).Error; err != nil {
		t.Fatalf("count withdrawals: %v", err)
	}
	if withdrawals != 1 {
		t.Fatalf("expected one withdrawal transaction, got %d", withdrawals)
	}
}
