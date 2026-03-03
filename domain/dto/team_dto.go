package dto

import (
	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/google/uuid"
)

type TeamCreate struct {
	InstiId       uuid.UUID `json:"insti_id"`
	Name          string    `json:"name"`
	Pelatih       string    `json:"pelatih"`
	LogoPath      string    `json:"logo_path"`
	RecLetterPath string    `json:"rec_letter_path"`
}

type TeamUpdate struct {
	Id            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Pelatih       string    `json:"pelatih"`
	LogoPath      string    `json:"logo_path"`
	RecLetterPath string    `json:"rec_letter_path"`
}

type TeamResponse struct {
	Id            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Pelatih       string    `json:"pelatih"`
	LogoPath      string    `json:"logo_path"`
	RecLetterPath string    `json:"rec_letter_path"`
	
	TeamMembers []TeamMemberResponse `json:"team_members"`
}

type TeamMemberCreate struct {
	TeamId uuid.UUID `json:"team_id"`
	FullName string `json:"full_name"`
	Role enums.TeamType `json:"role"`
	IdCardPath string `json:"id_card_path"`
}

type TeamMemberUpdate struct {
	Id uuid.UUID `json:"id"`
	FullName string `json:"full_name"`
	Role enums.TeamType `json:"role"`
	IdCardPath string `json:"id_card_path"`
}

type TeamMemberResponse struct {
	Id uuid.UUID `json:"id"`
	FullName string `json:"full_name"`
	Role enums.TeamType `json:"role"`
	IdCardPath string `json:"id_card_path"`
}	