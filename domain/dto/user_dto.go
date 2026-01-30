package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserRegister struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6,max=40,validpassword"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=6,max=40,validpassword"`
}

type UserParam struct {
	ID                  uuid.UUID `json:"id"`
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
	Token string `json:"token"`
}

type UserResponse struct {
	Id                  uuid.UUID          `json:"id"`
	Email               string             `json:"email"`
	Password            string             `json:"password"`
	EmailVerifiedToken  string             `json:"email_verified_token"`
	ForgotPasswordToken string             `json:"forgot_password_token"`
	EmailIsVerified     bool               `json:"email_is_verified"`
	ExpiredToken        time.Time          `json:"-"`
	ExpiredTokenForgot  time.Time          `json:"-"`
}

type UserUpdate struct {
	Password            string    `json:"password" binding:"omitempty,validpassword"`
	EmailIsVerified     bool      `json:"-"`
	EmailVerifiedToken  string    `json:"-"`
	ForgotPasswordToken string    `json:"-"`
	ExpiredToken        time.Time `json:"-"`
	ExpiredTokenForgot  time.Time `json:"-"`
}