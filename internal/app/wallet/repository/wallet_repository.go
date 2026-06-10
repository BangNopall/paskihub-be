package repository

import (
	"context"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) contracts.WalletRepository {
	return &walletRepository{
		db: db,
	}
}

func (r *walletRepository) CreateWallet(ctx context.Context, wallet *entity.Wallet) error {
	err := r.db.WithContext(ctx).Create(wallet).Error
	return err
}

func (r *walletRepository) GetWalletByEventId(ctx context.Context, eventId uuid.UUID) (*entity.Wallet, error) {
	var wallet entity.Wallet
	err := r.db.WithContext(ctx).Where("event_id = ?", eventId).First(&wallet).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepository) CreateTransaction(ctx context.Context, transaction *entity.WalletTransaction) error {
	err := r.db.WithContext(ctx).Create(transaction).Error
	return err
}

func (r *walletRepository) ApproveTransaction(ctx context.Context, transactionId uuid.UUID) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var transaction entity.WalletTransaction
		// Fetch transaction with lock for update
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", transactionId).First(&transaction).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return domain.ErrNotFound
			}
			return err
		}

		if transaction.Status != enums.Pending {
			return domain.ErrBadRequest // Status is not PENDING
		}

		// Update transaction status
		if err := tx.Model(&transaction).Update("status", enums.Approve).Error; err != nil {
			return err
		}

		// Fetch wallet with lock for update
		var wallet entity.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", transaction.WalletId).First(&wallet).Error; err != nil {
			return err
		}

		// Increase saldo and coin balance
		newSaldo := wallet.Saldo + transaction.Amount
		newCoinBalance := wallet.CoinBalance
		if transaction.CoinRateSnapshot > 0 {
			newCoinBalance += transaction.Amount / transaction.CoinRateSnapshot
		}

		if err := tx.Model(&wallet).Updates(map[string]interface{}{
			"saldo":        newSaldo,
			"coin_balance": newCoinBalance,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func (r *walletRepository) RejectTransaction(ctx context.Context, transactionId uuid.UUID, rejectionReason string) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var transaction entity.WalletTransaction
		// Fetch transaction with lock for update
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", transactionId).First(&transaction).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return domain.ErrNotFound
			}
			return err
		}

		if transaction.Status != enums.Pending {
			return domain.ErrBadRequest // Status is not PENDING
		}

		// Update transaction status and rejection reason
		updateData := map[string]interface{}{
			"status":           enums.TSRejected,
			"rejection_reason": rejectionReason,
		}

		if err := tx.Model(&transaction).Updates(updateData).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func (r *walletRepository) GetTransactionLogs(ctx context.Context, walletId uuid.UUID) ([]entity.WalletTransaction, error) {
	var transactions []entity.WalletTransaction
	err := r.db.WithContext(ctx).Where("wallet_id = ?", walletId).Order("created_at desc").Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *walletRepository) GetAllTransactions(ctx context.Context, status string, limit, offset int) ([]entity.WalletTransaction, int64, error) {
	var transactions []entity.WalletTransaction
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.WalletTransaction{}).
		Preload("Wallet").Preload("Wallet.Event")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at desc").Limit(limit).Offset(offset).Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *walletRepository) GetTransactionById(ctx context.Context, transactionId uuid.UUID) (*entity.WalletTransaction, error) {
	var transaction entity.WalletTransaction
	err := r.db.WithContext(ctx).Where("id = ?", transactionId).First(&transaction).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *walletRepository) CountTransactionsByStatus(ctx context.Context, walletId uuid.UUID, status string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.WalletTransaction{}).
		Where("wallet_id = ? AND status = ?", walletId, status).
		Count(&count).Error
	return count, err
}
