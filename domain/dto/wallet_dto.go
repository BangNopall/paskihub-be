package dto

import (
	"time"

	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/google/uuid"
)

type TopUpRequest struct {
	Amount     float64 `json:"amount" form:"amount" validate:"required,min=50000"`
	CouponCode string  `json:"coupon_code,omitempty" form:"coupon_code"`
}

type WalletResponse struct {
	Id        uuid.UUID `json:"id"`
	EventId   uuid.UUID `json:"event_id"`
	Saldo     float64   `json:"saldo"`
	SaldoKoin float64   `json:"saldo_koin"` // Saldo / 1000
}

type WalletTransactionResponse struct {
	Id         uuid.UUID               `json:"id"`
	WalletId   uuid.UUID               `json:"wallet_id"`
	Type       enums.WalletType        `json:"type"`
	Amount     float64                 `json:"amount"`
	AmountKoin float64                 `json:"amount_koin"` // Amount / 1000
	ProofPath  string                  `json:"proof_path"`
	Status     enums.TransactionStatus `json:"status"`
	CreatedAt  time.Time               `json:"created_at"`
	UpdatedAt  time.Time               `json:"updated_at"`
}

func WalletEntityToResponse(wallet *entity.Wallet) *WalletResponse {
	return &WalletResponse{
		Id:        wallet.Id,
		EventId:   wallet.EventId,
		Saldo:     wallet.Saldo,
		SaldoKoin: wallet.Saldo / 1000,
	}
}

func WalletTransactionEntityToResponse(transaction *entity.WalletTransaction) *WalletTransactionResponse {
	return &WalletTransactionResponse{
		Id:         transaction.Id,
		WalletId:   transaction.WalletId,
		Type:       transaction.Type,
		Amount:     transaction.Amount,
		AmountKoin: transaction.Amount / 1000,
		ProofPath:  transaction.ProofPath,
		Status:     transaction.Status,
		CreatedAt:  transaction.CreatedAt,
		UpdatedAt:  transaction.UpdatedAt,
	}
}
