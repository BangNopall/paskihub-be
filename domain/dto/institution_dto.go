package dto

import (
	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/google/uuid"
)

type InstitutionCreate struct {
	UserId  uuid.UUID             `json:"user_id"`
	Name    string                `json:"name"`
	Address string                `json:"address"`
	Type    enums.InstitutionType `json:"type"`
	NamePj  string                `json:"name_pj"`
	NoWaPj  string                `json:"no_wa_pj"`
}

type InstitutionUpdate struct {
	Id      uuid.UUID             `json:"id"`
	UserId  uuid.UUID             `json:"user_id"`
	Name    string                `json:"name"`
	Address string                `json:"address"`
	Type    enums.InstitutionType `json:"type"`
	NamePj  string                `json:"name_pj"`
	NoWaPj  string                `json:"no_wa_pj"`
}

type InstitutionResponse struct {
	Id      uuid.UUID             `json:"id"`
	Name    string                `json:"name"`
	Address string                `json:"address"`
	Type    enums.InstitutionType `json:"type"`
	NamePj  string                `json:"name_pj"`
	NoWaPj  string                `json:"no_wa_pj"`

	Users []UserResponse `json:"users"`
	Teams []TeamResponse `json:"teams"`
}
