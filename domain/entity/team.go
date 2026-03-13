package entity

import (
	"time"

	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/google/uuid"
)

type Team struct {
	Id            uuid.UUID `json:"id" gorm:"type:uuid;primarykey;"`
	InstiId       uuid.UUID `json:"insti_id" gorm:"type:uuid;index:idx_team_insti_id;"`
	Name          string    `json:"name" gorm:"type:varchar(255);"`
	LogoPath      string    `json:"logo_path" gorm:"type:varchar(255);"`
	Pelatih       string    `json:"pelatih" gorm:"type:varchar(255);"`
	RecLetterPath string    `json:"rec_letter_path" gorm:"type:varchar(255);"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;"`

	Institution Institution `json:"institution" gorm:"foreignKey:InstiId;references:Id;"`

	// Relation has many
	TeamMembers   []TeamMember   `json:"team_members" gorm:"foreignKey:TeamId;references:Id;"`
	Registrations []Registration `json:"registrations" gorm:"foreignKey:TeamId;references:Id;"`
}

type TeamMember struct {
	Id         uuid.UUID      `json:"id" gorm:"type:uuid;primarykey;"`
	TeamId     uuid.UUID      `json:"team_id" gorm:"type:uuid;index:idx_team_member_team_id;"`
	FullName   string         `json:"full_name" gorm:"type:varchar(255);"`
	Role       enums.TeamType `json:"role" gorm:"type:team_type;"`
	IdCardPath string         `json:"id_card_path" gorm:"type:varchar(255);"`
	PhotoPath  string         `json:"photo_path" gorm:"type:varchar(255);"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;"`

	Team Team `json:"team" gorm:"foreignKey:TeamId;references:Id;"`
}
