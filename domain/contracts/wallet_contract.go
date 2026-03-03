package contracts

import (
	"context"
	"mime/multipart"

	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/google/uuid"
)

type WalletRepository interface {
	CreateWallet(ctx context.Context, wallet *entity.Wallet) error
	GetWalletByEventId(ctx context.Context, eventId uuid.UUID) (*entity.Wallet, error)
	CreateTransaction(ctx context.Context, transaction *entity.WalletTransaction) error
	ApproveTransaction(ctx context.Context, transactionId uuid.UUID) error
	RejectTransaction(ctx context.Context, transactionId uuid.UUID) error
	GetTransactionLogs(ctx context.Context, walletId uuid.UUID) ([]entity.WalletTransaction, error)
	GetTransactionById(ctx context.Context, transactionId uuid.UUID) (*entity.WalletTransaction, error)
}

type WalletService interface {
	GetWalletInfo(ctx context.Context, eventId string, userId string) (*dto.WalletResponse, error)
	GetTransactionLogs(ctx context.Context, eventId string, userId string) ([]dto.WalletTransactionResponse, error)
	RequestTopUp(ctx context.Context, eventId string, userId string, req *dto.TopUpRequest, proofFile *multipart.FileHeader) error
	ApproveTopUp(ctx context.Context, transactionId string, adminUserId string) error
	RejectTopUp(ctx context.Context, transactionId string, adminUserId string) error
}
