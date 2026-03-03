package entity

import (
	"time"

	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/google/uuid"
)

type Wallet struct {
	Id        uuid.UUID `json:"id" gorm:"type:uuid;primarykey;"`
	EventId   uuid.UUID `json:"event_id" gorm:"type:uuid;index:idx_wallet_event_id;"`
	Saldo     float64   `json:"saldo" gorm:"type:decimal(15,2);default:0"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;"`

	WalletTransaction []WalletTransaction `json:"wallet_transaction" gorm:"foreignKey:WalletId;references:Id;"`
}

type WalletTransaction struct {
	Id        uuid.UUID               `json:"id" gorm:"type:uuid;primarykey;"`
	WalletId  uuid.UUID               `json:"wallet_id" gorm:"type:uuid;index:idx_wallet_transaction_wallet_id;"`
	Type      enums.WalletType        `json:"type" gorm:"type:wallet_type;"`
	Amount    float64                 `json:"amount" gorm:"type:decimal(15,2);"`
	ProofPath string                  `json:"proof_path" gorm:"type:varchar(255);default:null;"`
	Status    enums.TransactionStatus `json:"status" gorm:"type:transaction_status;"`
	CreatedAt time.Time               `json:"created_at" gorm:"autoCreateTime;"`
	UpdatedAt time.Time               `json:"updated_at" gorm:"autoUpdateTime;"`

	Wallet Wallet `json:"wallet" gorm:"foreignKey:WalletId;references:Id;"`
}
