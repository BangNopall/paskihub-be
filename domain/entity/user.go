package entity

import (
	"time"

	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/google/uuid"
)

type User struct {
	Id                  uuid.UUID  `json:"id" gorm:"type:uuid;primarykey;"`
	Email               string     `json:"email" gorm:"type:varchar(255);unique;not null;"`
	Role                enums.Role `json:"role" gorm:"type:role;"`
	Password            string     `json:"password" gorm:"type:varchar(255);not null;"`
	EmailVerifiedToken  string     `json:"email_verified_token" gorm:"type:varchar(100)"`
	ForgotPasswordToken string     `json:"forgot_password_token" gorm:"type:varchar(100)"`
	EmailIsVerified     bool       `json:"email_is_verified" gorm:"type:bool"`
	ExpiredToken        time.Time  `json:"-" gorm:"type:timestamp"`
	ExpiredTokenForgot  time.Time  `json:"-" gorm:"type:timestamp"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;"`

	// Relation has many
	Events       []Event       `json:"events" gorm:"foreignKey:UserId;references:Id;"`
	Institutions []Institution `json:"institutions" gorm:"foreignKey:UserId;references:Id;"`
}
