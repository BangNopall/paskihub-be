package dto

import (
	"time"

	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/google/uuid"
)

type UserRegister struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6,max=40,validpassword"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=6,max=40,validpassword"`
}

type UserParam struct {
	Id                  uuid.UUID `json:"id"`
	Email               string    `json:"email"`
	ForgotPasswordToken string    `json:"forgot_password_token"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserForgotPassword struct {
	Email string `json:"email" binding:"required,email"`
}

type UserResetPassword struct {
	Password        string `json:"password" binding:"required,min=6,max=20"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=6,max=20"`
}

type UserLoginResponse struct {
	Id       string  `json:"id"`
	Email    string  `json:"email"`
	Role     string  `json:"role"`
	ParentId *string `json:"parent_id,omitempty"`
	Token    string  `json:"token"`
}

type UserLoginInfoRes struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UserResponse struct {
	Id              uuid.UUID            `json:"id"`
	Email           string               `json:"email"`
	Role            string               `json:"role"`
	EmailIsVerified bool                 `json:"email_is_verified"`
	IsBanned        bool                 `json:"is_banned"`
	LastLoginAt     *time.Time           `json:"last_login_at"`
	CreatedAt       time.Time            `json:"created_at"`
	InstitutionName string               `json:"name"` // Used to hold dynamic name (Event/Institution/Admin)
	Event           *EventResponse       `json:"event,omitempty"`
	Institution     *InstitutionResponse `json:"institution,omitempty"`
}

type UserUpdate struct {
	Password            string     `json:"password" binding:"omitempty,validpassword"`
	LastLoginAt         *time.Time `json:"last_login_at"`
	EmailIsVerified     bool       `json:"-"`
	EmailVerifiedToken  string     `json:"-"`
	ForgotPasswordToken string     `json:"-"`
	ExpiredToken        time.Time  `json:"-"`
	ExpiredTokenForgot  time.Time  `json:"-"`
}

type UserPaginationResponse struct {
	Users      []UserResponse     `json:"users"`
	Pagination PaginationResponse `json:"pagination"`
}

// Admin Detail Response Structs
type AdminUserDetailResponse struct {
	Id              uuid.UUID               `json:"id"`
	Name            string                  `json:"name"`
	Email           string                  `json:"email"`
	Role            string                  `json:"role"`
	Status          string                  `json:"status"`
	EmailIsVerified bool                    `json:"email_is_verified"`
	IsBanned        bool                    `json:"is_banned"`
	LastLoginAt     *time.Time              `json:"last_login_at"`
	JoinedAt        time.Time               `json:"joined_at"`
	SchoolName      string                  `json:"school_name,omitempty"`
	Phone           string                  `json:"phone,omitempty"`
	Address         string                  `json:"address,omitempty"`
	Event           *EventResponse          `json:"event,omitempty"`
	Institution     *InstitutionResponse    `json:"institution,omitempty"`
	EOData          *AdminUserEODetail      `json:"eo_data,omitempty"`
	PesertaData     *AdminUserPesertaDetail `json:"peserta_data,omitempty"`
}

type AdminUserEODetail struct {
	Panitia []AdminStaffRes `json:"panitia"`
	Events  []AdminEventRes `json:"events"`
	Judges  []AdminJudgeRes `json:"judges"`
}

type AdminStaffRes struct {
	Name string `json:"name"`
	Role string `json:"role"` // Usually just "Staff"
}

type AdminEventLevelRes struct {
	Name     string `json:"name"`
	RegisFee string `json:"regis_fee"`
	DpFee    string `json:"dp_fee"`
}

type AdminEventRes struct {
	EventName   string               `json:"event_name"`
	CompeDate   time.Time            `json:"compe_date"`
	Location    string               `json:"location"`
	Status      string               `json:"status"`
	EventLevels []AdminEventLevelRes `json:"event_levels"`
}

type AdminJudgeRes struct {
	Name string `json:"name"`
}

type AdminUserPesertaDetail struct {
	Teams        []AdminTeamRes              `json:"teams"`
	EventHistory []AdminEventRegistrationRes `json:"event_history"`
}

type AdminTeamMemberRes struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type AdminTeamRes struct {
	TeamName     string               `json:"team_name"`
	Coach        string               `json:"coach"`
	MembersCount int                  `json:"members_count"`
	Members      []AdminTeamMemberRes `json:"members"`
}

type AdminEventRegistrationRes struct {
	EventName        string `json:"event_name"`
	EventLevelName   string `json:"event_level_name"`
	PaymentStatus    string `json:"payment_status"`
	AssessmentStatus string `json:"assessment_status"`
}

func UserEntityToResponse(user *entity.User) *UserResponse {
	var (
		name string = "-"
		evt  *EventResponse
		inst *InstitutionResponse
	)

	switch string(user.Role) {
	case "ADMIN":
		name = "Admin PaskiHub"
	case "ORGANIZER":
		if len(user.Events) > 0 {
			name = user.Events[0].Name
			evt = EventEntityToResponse(&user.Events[0])
			evt.User = nil // Remove redundancy as parent is the user itself
		}
	case "PESERTA":
		if len(user.Institutions) > 0 {
			name = user.Institutions[0].Name
			// Mapping to InstitutionResponse
			inst = &InstitutionResponse{
				Id:      user.Institutions[0].Id,
				Name:    user.Institutions[0].Name,
				Address: user.Institutions[0].Address,
				Type:    user.Institutions[0].InstitutionType,
				NamePj:  user.Institutions[0].NamePj,
				NoWaPj:  user.Institutions[0].NoWaPj,
			}
		}
	}

	return &UserResponse{
		Id:              user.Id,
		Email:           user.Email,
		Role:            string(user.Role),
		EmailIsVerified: user.EmailIsVerified,
		IsBanned:        user.IsBanned,
		LastLoginAt:     user.LastLoginAt,
		CreatedAt:       user.CreatedAt,
		InstitutionName: name,
		Event:           evt,
		Institution:     inst,
	}
}

type AdminCreateRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}
