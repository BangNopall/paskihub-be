package dto

import "mime/multipart"

type OpenEventResponse struct {
	Id          string                   `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	LogoPath    string                   `json:"logo_path"`
	PosterPath  string                   `json:"poster_path"`
	Levels      []OpenEventLevelResponse `json:"levels"`
}

type OpenEventLevelResponse struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	RegisFee string `json:"regis_fee"`
	DpFee    string `json:"dp_fee"`
}

type RegisterEventRequest struct {
	EventLevelId string                `form:"event_level_id" validate:"required"`
	TeamId       string                `form:"team_id" validate:"required"`
	PaymentType  string                `form:"payment_type" validate:"required"` // DP or LUNAS
	PaymentProof *multipart.FileHeader `swaggerignore:"true"`
}

type PelunasanEventRequest struct {
	PaymentProof *multipart.FileHeader `swaggerignore:"true"`
}

type ActiveEventResponse struct {
	RegistrationId string `json:"registration_id"`
	EventName      string `json:"event_name"`
	EventLogoPath  string `json:"event_logo_path"`
	TeamName       string `json:"team_name"`
	PaymentStatus  string `json:"payment_status"`
}
