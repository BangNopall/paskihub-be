package entity

import (
	"time"

	"github.com/google/uuid"
)

type Judge struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	EventID   uuid.UUID `gorm:"type:uuid;not null"`
	Name      string    `gorm:"type:varchar(255);not null"`
	UpdatedAt time.Time
	CreatedAt time.Time
}

type ScoreCategory struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	EventID       uuid.UUID `gorm:"type:uuid;not null"`
	EventLevelID  uuid.UUID `gorm:"type:uuid;not null"`
	Name          string    `gorm:"type:varchar(255);not null"`
	UpdatedAt     time.Time
	CreatedAt     time.Time
	SubCategories []ScoreSubCategory `gorm:"foreignKey:ScoreCategoriesID;references:ID"`
}

type ScoreSubCategory struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ScoreCategoriesID uuid.UUID `gorm:"type:uuid;not null"`
	Name              string    `gorm:"type:varchar(255);not null"`
	MaxScore          float64   `gorm:"not null;default:0"`
	UpdatedAt         time.Time
	CreatedAt         time.Time
	GradeRules        []GradeRule `gorm:"foreignKey:ScoreSubCategoryID;references:ID"`
}

type ViolationType struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	EventID      uuid.UUID `gorm:"type:uuid;not null"`
	EventLevelID uuid.UUID `gorm:"type:uuid;not null"`
	Name         string    `gorm:"type:varchar(255);not null"`
	Point        float64   `gorm:"not null"`
	UpdatedAt    time.Time
	CreatedAt    time.Time
}

type GradeRule struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	EventID           uuid.UUID `gorm:"type:uuid;not null"`
	EventLevelID      uuid.UUID `gorm:"type:uuid;not null"`
	ScoreSubCategoryID uuid.UUID `gorm:"type:uuid;default:null"`
	GradeName         string    `gorm:"type:varchar(50);not null"` // Hardcoded "Kurang", "Cukup", "Baik", "Sangat Baik"
	MinScore          float64   `gorm:"not null"`
	MaxScore          float64   `gorm:"not null"`
	UpdatedAt         time.Time
	CreatedAt         time.Time
}

type Score struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	RegisID       uuid.UUID `gorm:"type:uuid;not null;index:idx_score_unique,unique"`
	JudgesID      uuid.UUID `gorm:"type:uuid;not null;index:idx_score_unique,unique"`
	SubCategoryID uuid.UUID `gorm:"type:uuid;not null;index:idx_score_unique,unique"`
	ScoreValue    float64   `gorm:"not null"`
	Grade         string    `gorm:"type:varchar(50);not null"`
	UpdatedAt     time.Time
	CreatedAt     time.Time

	Registration     Registration     `json:"registration" gorm:"foreignKey:RegisID;references:Id;"`
	Judge            Judge            `json:"judge" gorm:"foreignKey:JudgesID;references:ID;"`
	ScoreSubCategory ScoreSubCategory `json:"score_sub_category" gorm:"foreignKey:SubCategoryID;references:ID;"`
}

type TeamViolation struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	RegisID         uuid.UUID `gorm:"type:uuid;not null"`
	JudgesID        uuid.UUID `gorm:"type:uuid;not null"`
	ViolationTypeID uuid.UUID `gorm:"type:uuid;not null"`
	UpdatedAt       time.Time
	CreatedAt       time.Time

	Registration  Registration  `json:"registration" gorm:"foreignKey:RegisID;references:Id;"`
	Judge         Judge         `json:"judge" gorm:"foreignKey:JudgesID;references:ID;"`
	ViolationType ViolationType `json:"violation_type" gorm:"foreignKey:ViolationTypeID;references:ID;"`
}

type EventAward struct {
	ID              uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	EventID         uuid.UUID       `gorm:"type:uuid;not null"`
	EventLevelID    uuid.UUID       `gorm:"type:uuid;not null"`
	Name            string          `gorm:"type:varchar(255);not null"`
	LimitRank       int             `gorm:"not null;default:1"`
	CreatedAt       time.Time
	UpdatedAt       time.Time

	// Relation
	Event           Event           `gorm:"foreignKey:EventID;references:Id"`
	EventLevel      EventLevel      `gorm:"foreignKey:EventLevelID;references:Id"`
	ScoreCategories []ScoreCategory `gorm:"many2many:event_award_score_categories;"`
}
