package dto

import (
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/google/uuid"
)

type RegistrationCreate struct {
	EventLevelId     uuid.UUID `json:"event_level_id"`
	TeamId           uuid.UUID `json:"team_id"`
	PaymentProofPath string    `json:"payment_proof_path"`
}

type RegistrationUpdate struct {
	Id               uuid.UUID                `json:"id"`
	EventLevelId     uuid.UUID                `json:"event_level_id"`
	TeamId           uuid.UUID                `json:"team_id"`
	PaymentStatus    enums.RegistrationStatus `json:"payment_status"`
	PaymentProofPath string                   `json:"payment_proof_path"`
	RejectionReason  string                   `json:"rejection_reason"`
	IsKick           bool                     `json:"is_kick"`
}

type RegistrationResponse struct {
	Id               uuid.UUID                `json:"id"`
	EventLevelId     uuid.UUID                `json:"event_level_id"`
	TeamId           uuid.UUID                `json:"team_id"`
	PaymentStatus    enums.RegistrationStatus `json:"payment_status"`
	PaymentProofPath string                   `json:"payment_proof_path"`
	RejectionReason  string                   `json:"rejection_reason"`
	IsKick           bool                     `json:"is_kick"`

	Team       TeamResponse       `json:"team"`
	EventLevel EventLevelResponse `json:"event_level"`
}

type RegistrationPaginationResponse struct {
	Registrations []RegistrationResponse `json:"registrations"`
	Pagination    PaginationResponse     `json:"pagination"`
}

func RegistrationEntityToResponse(registration *entity.Registration) *RegistrationResponse {
	return &RegistrationResponse{
		Id:               registration.Id,
		EventLevelId:     registration.EventLevelId,
		TeamId:           registration.TeamId,
		PaymentStatus:    registration.PaymentStatus,
		PaymentProofPath: registration.PaymentProofPath,
		RejectionReason:  registration.RejectionReason,
		IsKick:           registration.IsKick,
	}
}
