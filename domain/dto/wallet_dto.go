package dto

import (
	"time"

	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/google/uuid"
)

type TopUpRequest struct {
	Amount     float64 `json:"amount" form:"amount" validate:"required,min=1"`
	CouponCode string  `json:"coupon_code,omitempty" form:"coupon_code"`
}

type WalletResponse struct {
	Id                   uuid.UUID `json:"id"`
	EventId              uuid.UUID `json:"event_id"`
	Saldo                float64   `json:"saldo"`
	SaldoKoin            float64   `json:"saldo_koin"` // Saldo / CoinRate
	SuccessfulTopupCount int64     `json:"successful_topup_count"`
	PendingTopupCount    int64     `json:"pending_topup_count"`
}

type WalletTransactionResponse struct {
	Id              uuid.UUID               `json:"id"`
	WalletId        uuid.UUID               `json:"wallet_id"`
	Type            enums.WalletType        `json:"type"`
	Amount          float64                 `json:"amount"`
	AmountKoin      float64                 `json:"amount_koin"` // Amount / CoinRate
	ProofPath       string                  `json:"proof_path"`
	Status          enums.TransactionStatus `json:"status"`
	RejectionReason string                  `json:"rejection_reason"`
	CreatedAt       time.Time               `json:"created_at"`
	UpdatedAt       time.Time               `json:"updated_at"`
}

type AdminTransactionResponse struct {
	Id              uuid.UUID               `json:"id"`
	EOName          string                  `json:"eo_name"`
	Amount          float64                 `json:"amount"`
	AmountKoin      float64                 `json:"amount_koin"`
	Status          enums.TransactionStatus `json:"status"`
	ProofPath       string                  `json:"proof_path"`
	RejectionReason string                  `json:"rejection_reason"`
	CreatedAt       time.Time               `json:"created_at"`
}

type RejectTopUpRequest struct {
	RejectionReason string `json:"rejection_reason" validate:"required"`
}

type AdminTransactionPaginationResponse struct {
	Transactions []AdminTransactionResponse `json:"transactions"`
	Pagination   PaginationResponse         `json:"pagination"`
}

func WalletEntityToResponse(wallet *entity.Wallet, coinRate float64) *WalletResponse {
	return &WalletResponse{
		Id:        wallet.Id,
		EventId:   wallet.EventId,
		Saldo:     wallet.Saldo,
		SaldoKoin: wallet.Saldo / coinRate,
	}
}

func WalletTransactionEntityToResponse(transaction *entity.WalletTransaction, coinRate float64) *WalletTransactionResponse {
	return &WalletTransactionResponse{
		Id:              transaction.Id,
		WalletId:        transaction.WalletId,
		Type:            transaction.Type,
		Amount:          transaction.Amount,
		AmountKoin:      transaction.Amount / coinRate,
		ProofPath:       transaction.ProofPath,
		Status:          transaction.Status,
		RejectionReason: transaction.RejectionReason,
		CreatedAt:       transaction.CreatedAt,
		UpdatedAt:       transaction.UpdatedAt,
	}
}
