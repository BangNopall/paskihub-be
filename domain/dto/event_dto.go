package dto

import (
	"time"

	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/google/uuid"
)

type EventCreate struct {
	UserId     uuid.UUID `json:"user_id"`
	Name       string    `json:"name"`
	OpenDate   string    `json:"open_date"`
	CloseDate  string    `json:"close_date"`
	Location   string    `json:"location"`
	NamaPj     string    `json:"nama_pj"`
	NoWaPj     string    `json:"no_wa_pj"`
	BankName   string    `json:"bank_name"`
	BankNumber string    `json:"bank_number"`
}

type EventParams struct {
	Id     uuid.UUID
	Name   string
	Status string
}

type EventUpdate struct {
	UserId         uuid.UUID `json:"user_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	LogoPath       string    `json:"-"`
	PosterPath     string    `json:"-"`
	OpenDate       string    `json:"open_date"`
	CloseDate      string    `json:"close_date"`
	CompeDate      string    `json:"compe_date"`
	Organizer      string    `json:"organizer"`
	Location       string    `json:"location"`
	MinTeamMembers int       `json:"min_team_members"`
	MaxTeamMembers int       `json:"max_team_members"`
	Status         string    `json:"status"`
	NamePj         string    `json:"name_pj"`
	NoWaPj         string    `json:"no_wa_pj"`
	BankName       string    `json:"bank_name"`
	BankNumber     string    `json:"bank_number"`
	WaGroup        string    `json:"wa_group"`
}

type EventResponse struct {
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	LogoPath       string                 `json:"logo_path"`
	PosterPath     string                 `json:"poster_path"`
	OpenDate       time.Time              `json:"open_date"`
	CloseDate      time.Time              `json:"close_date"`
	CompeDate      time.Time              `json:"compe_date"`
	Organizer      string                 `json:"organizer"`
	Location       string                 `json:"location"`
	MinTeamMembers int                    `json:"min_team_members"`
	MaxTeamMembers int                    `json:"max_team_members"`
	Status         string                 `json:"status"`
	NamePj         string                 `json:"name_pj"`
	NoWaPj         string                 `json:"no_wa_pj"`
	BankName       string                 `json:"bank_name"`
	BankNumber     string                 `json:"bank_number"`
	WaGroup        string                 `json:"wa_group"`
	User           UserResponse           `json:"user"`
	EventLevels    []EventLevelResponse   `json:"event_levels"`
	Registrations  []RegistrationResponse `json:"registrations"`
}

type EventLevelCreate struct {
	EventId  uuid.UUID `json:"event_id"`
	Name     string    `json:"name"`
	RegisFee string    `json:"regis_fee"`
	DpFee    string    `json:"dp_fee"`
}

type EventLevelUpdate struct {
	Id       uuid.UUID `json:"id"`
	EventId  uuid.UUID `json:"event_id"`
	Name     string    `json:"name"`
	RegisFee string    `json:"regis_fee"`
	DpFee    string    `json:"dp_fee"`
}

type EventLevelResponse struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	RegisFee string    `json:"regis_fee"`
	DpFee    string    `json:"dp_fee"`
}

type EventPaginationResponse struct {
	Events     []EventResponse    `json:"events"`
	Pagination PaginationResponse `json:"pagination"`
}

func EventLevelEntityToResponse(eventLevel *entity.EventLevel) *EventLevelResponse {
	return &EventLevelResponse{
		Id:       eventLevel.Id,
		Name:     eventLevel.Name,
		RegisFee: eventLevel.RegisFee,
		DpFee:    eventLevel.DpFee,
	}
}

func EventEntityToResponse(event *entity.Event) *EventResponse {
	return &EventResponse{
		Name:           event.Name,
		Description:    event.Description,
		LogoPath:       event.LogoPath,
		PosterPath:     event.PosterPath,
		OpenDate:       event.OpenDate,
		CloseDate:      event.CloseDate,
		CompeDate:      event.CompeDate,
		Organizer:      event.Organizer,
		Location:       event.Location,
		MinTeamMembers: event.MinTeamMembers,
		MaxTeamMembers: event.MaxTeamMembers,
		Status:         event.Status,
		NamePj:         event.NamePj,
		NoWaPj:         event.NoWaPj,
		BankName:       event.BankName,
		BankNumber:     event.BankNumber,
		WaGroup:        event.WaGroup,
		User:           *UserEntityToResponse(&event.User),
		EventLevels: func() []EventLevelResponse {
			res := make([]EventLevelResponse, 0)
			for _, eventLevel := range event.EventLevels {
				res = append(res, *EventLevelEntityToResponse(&eventLevel))
			}
			return res
		}(),
		Registrations: func() []RegistrationResponse {
			res := make([]RegistrationResponse, 0)
			for _, registration := range event.Registrations {
				res = append(res, *RegistrationEntityToResponse(&registration))
			}
			return res
		}(),
	}
}
