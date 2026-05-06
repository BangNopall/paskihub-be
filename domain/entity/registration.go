package entity

import (
	"time"

	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/google/uuid"
)

type Registration struct {
	Id               uuid.UUID                `json:"id" gorm:"type:uuid;primarykey;"`
	TeamId           uuid.UUID                `json:"team_id" gorm:"type:uuid;index:idx_registration_team_id;"`
	EventLevelId     uuid.UUID                `json:"event_level_id" gorm:"type:uuid;index:idx_registration_event_level_id;"`
	PaymentType      string                   `json:"payment_type" gorm:"type:varchar(20);"` // DP or LUNAS
	PaymentStatus    enums.RegistrationStatus `json:"payment_status" gorm:"type:registration_status;"`
	AssessmentStatus enums.AssessmentStatus   `json:"assessment_status" gorm:"type:varchar(20);default:'PENDING';"`
	PaymentProofPath string                   `json:"payment_proof_path" gorm:"type:varchar(255);"`
	RejectionReason  string                   `json:"rejection_reason" gorm:"type:varchar(255);"`
	IsKick           bool                     `json:"is_kick" gorm:"type:bool;default:false;"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;"`

	Team       Team       `json:"team" gorm:"foreignKey:TeamId;references:Id;"`
	EventLevel EventLevel `json:"event_level" gorm:"foreignKey:EventLevelId;references:Id;"`
}
