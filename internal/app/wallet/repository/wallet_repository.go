package repository

import (
	"context"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", transactionId).First(&transaction).Error; err != nil {
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
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", transaction.WalletId).First(&wallet).Error; err != nil {
			return err
		}

		// Increase saldo
		newSaldo := wallet.Saldo + transaction.Amount
		if err := tx.Model(&wallet).Update("saldo", newSaldo).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func (r *walletRepository) RejectTransaction(ctx context.Context, transactionId uuid.UUID) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var transaction entity.WalletTransaction
		// Fetch transaction with lock for update
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", transactionId).First(&transaction).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return domain.ErrNotFound
			}
			return err
		}

		if transaction.Status != enums.Pending {
			return domain.ErrBadRequest // Status is not PENDING
		}

		// Update transaction status
		if err := tx.Model(&transaction).Update("status", enums.TSRejected).Error; err != nil {
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
