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
	ID            uuid.UUID          `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	EventID       uuid.UUID          `gorm:"type:uuid;not null"`
	Name          string             `gorm:"type:varchar(255);not null"`
	UpdatedAt     time.Time
	CreatedAt     time.Time
	SubCategories []ScoreSubCategory `gorm:"foreignKey:ScoreCategoriesID;references:ID"`
}

type ScoreSubCategory struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ScoreCategoriesID uuid.UUID `gorm:"type:uuid;not null"`
	Name              string    `gorm:"type:varchar(255);not null"`
	UpdatedAt         time.Time
	CreatedAt         time.Time
}

type ViolationType struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	EventID   uuid.UUID `gorm:"type:uuid;not null"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Point     float64   `gorm:"not null"`
	UpdatedAt time.Time
	CreatedAt time.Time
}

type GradeRule struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	EventID   uuid.UUID `gorm:"type:uuid;not null"`
	GradeName string    `gorm:"type:varchar(50);not null"` // Hardcoded "Kurang", "Cukup", "Baik", "Sangat Baik"
	MinScore  float64   `gorm:"not null"`
	MaxScore  float64   `gorm:"not null"`
	UpdatedAt time.Time
	CreatedAt time.Time
}

type Score struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	RegisID       uuid.UUID `gorm:"type:uuid;not null"`
	JudgesID      uuid.UUID `gorm:"type:uuid;not null"`
	SubCategoryID uuid.UUID `gorm:"type:uuid;not null"`
	ScoreValue    float64   `gorm:"not null"`
	Grade         string    `gorm:"type:varchar(50);not null"`
	UpdatedAt     time.Time
	CreatedAt     time.Time
}
