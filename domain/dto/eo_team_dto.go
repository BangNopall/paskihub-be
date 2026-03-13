package dto

import (
	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/google/uuid"
)

type EOTeamRejectReq struct {
	RejectionReason string `json:"rejection_reason" validate:"required"`
}

type EOTeamApproveReq struct {
	PaymentStatus enums.RegistrationStatus `json:"payment_status" validate:"required,oneof=DP_PAID FULL_PAID"`
}

type EOTeamListRes struct {
	RegistrationId uuid.UUID                `json:"registration_id"`
	TeamId         uuid.UUID                `json:"team_id"`
	LogoPath       string                   `json:"logo_path"`
	TeamName       string                   `json:"team_name"`
	Institution    string                   `json:"institution"`
	EventLevel     string                   `json:"event_level"`
	PaymentStatus  enums.RegistrationStatus `json:"payment_status"`
}

type EOTeamDetailRes struct {
	RegistrationId   uuid.UUID                `json:"registration_id"`
	TeamId           uuid.UUID                `json:"team_id"`
	TeamName         string                   `json:"team_name"`
	LogoPath         string                   `json:"logo_path"`
	Pelatih          string                   `json:"pelatih"`
	RecLetterPath    string                   `json:"rec_letter_path"`
	Institution      string                   `json:"institution"`
	EventLevel       string                   `json:"event_level"`
	PaymentStatus    enums.RegistrationStatus `json:"payment_status"`
	PaymentProofPath string                   `json:"payment_proof_path"`
	RejectionReason  string                   `json:"rejection_reason"`
	IsKick           bool                     `json:"is_kick"`
	Members          []EOTeamMemberRes        `json:"members"`
}

type EOTeamMemberRes struct {
	Id         uuid.UUID      `json:"id"`
	FullName   string         `json:"full_name"`
	Role       enums.TeamType `json:"role"`
	IdCardPath string         `json:"id_card_path"`
}
