package entity

import (
	"time"

	"github.com/google/uuid"
)

type Jury struct {
	Id        uuid.UUID `json:"id" gorm:"type:uuid;primarykey;"`
	UserId    uuid.UUID `json:"user_id" gorm:"type:uuid;index:idx_jury_user_id;"` // Optional: if juries are also users
	Name      string    `json:"name" gorm:"type:varchar(255);"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;"`

	User User `json:"user" gorm:"foreignKey:UserId;references:Id;"`
}
