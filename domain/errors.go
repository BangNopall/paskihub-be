package domain

import (
	"errors"
	"net/http"
)

var (
	ErrAnnouncementNotFound    = errors.New("announcement not found")
	ErrCompetitionNotFound     = errors.New("competition not found")
	ErrTeamNotFound            = errors.New("team not found")
	ErrUserNotTeamLeader       = errors.New("user not a team leader")
	ErrVoucherNotFound         = errors.New("invalid voucher")
	ErrVoucherAlreadyRedeemed  = errors.New("voucher already redeemed")
	ErrInvalidCompeTeamID      = errors.New("invalid compe or team id")
	ErrIllegalEntry            = errors.New("illegal entry")
	ErrUserAlreadyRegistered   = errors.New("user already registered")
	ErrInternalServer          = errors.New("internal server error")
	ErrEmailRegistered         = errors.New("email already registered")
	ErrCheckEmail              = errors.New("please check your email to verify")
	ErrInvalidToken            = errors.New("invalid token")
	ErrWrongEmailOrPassword    = errors.New("invalid username/password")
	ErrFileTooBig              = errors.New("file size too big")
	ErrNotFound                = errors.New("item not found")
	ErrTimeout                 = errors.New("operation timed out")
	ErrDuplicateEntry          = errors.New("data already exists")
	ErrUserNotFound            = errors.New("user not found")
	ErrUniversityNotFound      = errors.New("university not found")
	ErrAlreadyAttend           = errors.New("user already attended")
	ErrForbiddenUpdate         = errors.New("forbidden to update")
	ErrInvalidEnumInput        = errors.New("invalid enum input")
	ErrMissingAttribute        = errors.New("missing required attribute")
	ErrConfirmPasswordNotMatch = errors.New("password and confirm password doesn't match")
	ErrTeamFull                = errors.New("team is full")
	ErrInvalidProofType        = errors.New("invalid proof type")
	ErrBadRequest              = errors.New("bad data request")
	ErrInvalidRole             = errors.New("invalid user data")
	ErrForbidden               = errors.New("forbidden access")
	ErrAccountBanned           = errors.New("account is banned")
	ErrRegistrationNotFound        = errors.New("registration not found for this event")
	ErrInsufficientBalance         = errors.New("insufficient wallet balance for approval fee")
	ErrTeamNotBelongToInstitution  = errors.New("team does not belong to your institution")
	ErrInstitutionCategoryMismatch = errors.New("team institution category does not match event category")
	ErrInvalidTeamMembersCount     = errors.New("invalid team members count for this event")
	ErrAlreadyRegisteredForEvent   = errors.New("team is already registered for this event")
	ErrAlreadyFullyPaid            = errors.New("registration is already fully paid")
)

func GetCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if errors.Is(err, ErrInternalServer) {
		return http.StatusInternalServerError
	}
	if errors.Is(err, ErrNotFound) || errors.Is(err, ErrVoucherNotFound) || errors.Is(err, ErrUserNotFound) {
		return http.StatusNotFound
	}
	if errors.Is(err, ErrInvalidEnumInput) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrUserNotTeamLeader) || errors.Is(err, ErrAnnouncementNotFound) ||
		errors.Is(err, ErrCompetitionNotFound) || errors.Is(err, ErrTeamNotFound) ||
		errors.Is(err, ErrIllegalEntry) || errors.Is(err, ErrMissingAttribute) ||
		errors.Is(err, ErrTeamFull) || errors.Is(err, ErrInvalidCompeTeamID) ||
		errors.Is(err, ErrInvalidProofType) || errors.Is(err, ErrUniversityNotFound) ||
		errors.Is(err, ErrTeamNotBelongToInstitution) || errors.Is(err, ErrInstitutionCategoryMismatch) ||
		errors.Is(err, ErrInvalidTeamMembersCount) || errors.Is(err, ErrAlreadyFullyPaid) ||
		errors.Is(err, ErrBadRequest) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrTimeout) {
		return http.StatusRequestTimeout
	}
	if errors.Is(err, ErrDuplicateEntry) || errors.Is(err, ErrUserAlreadyRegistered) ||
		errors.Is(err, ErrEmailRegistered) || errors.Is(err, ErrVoucherAlreadyRedeemed) ||
		errors.Is(err, ErrAlreadyRegisteredForEvent) {
		return http.StatusConflict
	}
	if errors.Is(err, ErrAlreadyAttend) {
		return http.StatusUnprocessableEntity
	}
	if errors.Is(err, ErrCheckEmail) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrInvalidToken) || errors.Is(err, ErrWrongEmailOrPassword) {
		return http.StatusUnauthorized
	}
	if errors.Is(err, ErrFileTooBig) {
		return http.StatusRequestEntityTooLarge
	}
	if errors.Is(err, ErrForbiddenUpdate) || errors.Is(err, ErrForbidden) || errors.Is(err, ErrAccountBanned) {
		return http.StatusForbidden
	}
	if errors.Is(err, ErrConfirmPasswordNotMatch) || errors.Is(err, ErrInvalidRole) {
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}
