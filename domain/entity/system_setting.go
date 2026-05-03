package entity

import (
	"time"

	"github.com/google/uuid"
)

type SystemSetting struct {
	Id            uuid.UUID `json:"id" gorm:"type:uuid;primarykey;"`
	CoinRate      float64   `json:"coin_rate" gorm:"type:decimal(15,2);default:1000"`
	ApprovalFee   float64   `json:"approval_fee" gorm:"type:decimal(15,2);default:1"`
	BankName      string    `json:"bank_name" gorm:"type:varchar(255)"`
	AccountNumber string    `json:"account_number" gorm:"type:varchar(255)"`
	AccountName   string    `json:"account_name" gorm:"type:varchar(255)"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime;"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime;"`
}
