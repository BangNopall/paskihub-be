package entity

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	Id             uuid.UUID `json:"id" gorm:"type:uuid;primarykey;"`
	UserId         uuid.UUID `json:"user_id" gorm:"type:uuid;index:idx_event_user_id;"`
	Name           string    `json:"name" gorm:"type:varchar(255);"`
	Description    string    `json:"description" gorm:"type:text;default:null;"`
	LogoPath       string    `json:"logo_path" gorm:"type:varchar(255);default:null;"`
	PosterPath     string    `json:"poster_path" gorm:"type:varchar(255);default:null;"`
	OpenDate       time.Time `json:"open_date" gorm:"type:timestamp;default:null;"`
	CloseDate      time.Time `json:"close_date" gorm:"type:timestamp;default:null;"`
	CompeDate      time.Time `json:"compe_date" gorm:"type:timestamp;default:null;"`
	Organizer      string    `json:"organizer" gorm:"type:varchar(255);default:null;"`
	Location       string    `json:"location" gorm:"type:varchar(255);default:null;"`
	MinTeamMembers int       `json:"min_team_members" gorm:"type:int;default:null;"`
	MaxTeamMembers int       `json:"max_team_members" gorm:"type:int;default:null;"`
	Status         string    `json:"status" gorm:"type:event_status;default:DRAFT;"`
	NamePj         string    `json:"name_pj" gorm:"type:varchar(255);default:null;"`
	NoWaPj         string    `json:"no_wa_pj" gorm:"type:varchar(255);default:null;"`
	BankName       string    `json:"bank_name" gorm:"type:varchar(255);default:null;"`
	BankNumber     string    `json:"bank_number" gorm:"type:varchar(255);default:null;"`
	WaGroup        string    `json:"wa_group" gorm:"type:varchar(255);default:null;"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime;"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime;"`

	User          User           `json:"user" gorm:"foreignKey:UserId;references:Id;"`
	
	// Relation has many
	EventLevels   []EventLevel   `json:"event_levels" gorm:"foreignKey:EventId;references:Id;"`
	Registrations []Registration `json:"registrations" gorm:"-"`
}

type EventLevel struct {
	Id        uuid.UUID `json:"id" gorm:"type:uuid;primarykey"`
	EventId   uuid.UUID `json:"event_id" gorm:"type:uuid;index:idx_event_level_id"`
	Name      string    `json:"name" gorm:"type:varchar(255);"`
	RegisFee         string    `json:"regis_fee" gorm:"type:varchar(255)"`
	DpFee            string    `json:"dp_fee" gorm:"type:varchar(255)"`
	IsScorePublished bool      `json:"is_score_published" gorm:"type:bool;default:false"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime;"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;"`

	Event Event `json:"event" gorm:"foreignKey:EventId;references:Id;"`

	// Relation has many
	Registrations []Registration `json:"registrations" gorm:"foreignKey:EventLevelId;references:Id;"`
}
