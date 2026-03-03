package entity

import (
	"time"

	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/google/uuid"
)

type Institution struct {
	Id              uuid.UUID             `json:"id" gorm:"type:uuid;primarykey;"`
	UserId          uuid.UUID             `json:"user_id" gorm:"type:uuid;index:idx_institution_user_id;"`
	Name            string                `json:"name" gorm:"type:varchar(255);"`
	Address         string                `json:"address" gorm:"type:text;"`
	InstitutionType enums.InstitutionType `json:"institution_type" gorm:"type:institution_type;"`
	NamePj          string                `json:"name_pj" gorm:"type:varchar(255);"`
	NoWaPj          string                `json:"no_wa_pj" gorm:"type:varchar(255);"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;"`

	User User `json:"user" gorm:"foreignKey:UserId;references:Id;"`

	// Relation has many
	Teams []Team `json:"teams" gorm:"foreignKey:InstiId;references:Id;"`
}
