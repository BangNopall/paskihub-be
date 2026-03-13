package dto

import "mime/multipart"

type CreateTeamRequest struct {
	Name             string
	PelatihName      string
	LogoTeam         *multipart.FileHeader
	SuratRekomendasi *multipart.FileHeader
	Members          []TeamMemberRequest
}

type TeamMemberRequest struct {
	FullName string
	Role     string
	IdCard   *multipart.FileHeader
	Photo    *multipart.FileHeader
}

type ParticipantTeamResponse struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	LogoPath      string `json:"logo_path"`
	PaymentStatus string `json:"payment_status"`
}

type TeamDetailResponse struct {
	Id             string                          `json:"id"`
	Name           string                          `json:"name"`
	LogoPath       string                          `json:"logo_path"`
	Pelatih        string                          `json:"pelatih"`
	RecLetterPath  string                          `json:"rec_letter_path"`
	MembersGrouped map[string][]ParticipantTeamMemberResponse `json:"members_grouped"`
}

type ParticipantTeamMemberResponse struct {
	Id         string `json:"id"`
	FullName   string `json:"full_name"`
	Role       string `json:"role"`
	IdCardPath string `json:"id_card_path"`
	PhotoPath  string `json:"photo_path"`
}
