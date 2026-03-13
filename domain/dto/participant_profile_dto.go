package dto

type ParticipantProfileResponse struct {
	Email       string                      `json:"email"`
	Institution *InstitutionProfileResponse `json:"institution"`
}

type InstitutionProfileResponse struct {
	Name            string `json:"name"`
	Address         string `json:"address"`
	InstitutionType string `json:"institution_type"`
	NamePj          string `json:"name_pj"`
	NoWaPj          string `json:"no_wa_pj"`
}

type UpdateInstitutionRequest struct {
	Name            string `json:"name" validate:"required"`
	Address         string `json:"address" validate:"required"`
	InstitutionType string `json:"institution_type" validate:"required"`
	NamePj          string `json:"name_pj" validate:"required"`
	NoWaPj          string `json:"no_wa_pj" validate:"required"`
}

type UpdatePasswordRequest struct {
	OldPassword     string `json:"old_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}
