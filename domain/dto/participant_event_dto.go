package dto

import "mime/multipart"

type OpenEventResponse struct {
	Id             string                   `json:"id"`
	Name           string                   `json:"name"`
	Description    string                   `json:"description"`
	LogoPath       string                   `json:"logo_path"`
	PosterPath     string                   `json:"poster_path"`
	Organizer      string                   `json:"organizer"`
	Status         string                   `json:"status"`
	OpenDate       string                   `json:"open_date"`
	CloseDate      string                   `json:"close_date"`
	Location       string                   `json:"location"`
	MinTeamMembers int                      `json:"min_team_members"`
	MaxTeamMembers int                      `json:"max_team_members"`
	Levels         []OpenEventLevelResponse `json:"levels"`
	BankName       string                   `json:"bank_name"`
	BankNumber     string                   `json:"bank_number"`
	NamePj         string                   `json:"name_pj"`
	NoWaPj         string                   `json:"no_wa_pj"`
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
	RegistrationId  string `json:"registration_id"`
	EventName       string `json:"event_name"`
	EventLogoPath   string `json:"event_logo_path"`
	TeamName        string `json:"team_name"`
	PaymentStatus   string `json:"payment_status"`
	PaymentType     string `json:"payment_type"`
	RejectionReason string `json:"rejection_reason"`
	IsKick          bool   `json:"is_kick"`
}

type RegistrationDetailResponse struct {
	EventLevelId string `json:"event_level_id"`
	Event        struct {
		Id          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Date        string `json:"date"`
		Location    string `json:"location"`
		Price       string `json:"price"`
		TargetDate  string `json:"target_date"`
		LogoUrl     string `json:"logo_url"`
	} `json:"event"`
	Team struct {
		Id            string `json:"id"`
		Name          string `json:"name"`
		LogoUrl       string `json:"logo_url"`
		OfficialCount int    `json:"official_count"`
		PasukanCount  int    `json:"pasukan_count"`
	} `json:"team"`
	Payment struct {
		Status          string `json:"status"`
		AmountPaid      string `json:"amount_paid"`
		TotalAmount     string `json:"total_amount"`
		RemainingAmount string `json:"remaining_amount"`
		ProofUrl        string `json:"proof_url"`
	} `json:"payment"`
}
