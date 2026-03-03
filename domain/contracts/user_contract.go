package contracts

import (
	// "github.com/gofiber/fiber/v2"
	"context"
	// "mime/multipart"

	"github.com/google/uuid"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
)

type UserRepository interface {
	CreateUser(user *entity.User) error
	// FetchAllByConditionAndRelation(
	// 	condition string,
	// 	args []interface{},
	// 	joins []string,
	// 	pageParam *dto.PaginationRequest,
	// 	preload ...string,
	// ) ([]entity.User, dto.PaginationResponse, error)
	FindUser(user *entity.User, userParam *dto.UserParam, relations ...string) error
	UpdateUser(updateUser *dto.UserUpdate, userId uuid.UUID) error
	DeleteUnverifiedUser() error
}

type UserService interface {
	Register(ctx context.Context, role string, user dto.UserRegister) error
	VerifyEmail(ctx context.Context, email string, emailVerPass string) error
	Login(ctx context.Context, user dto.UserLogin) (dto.UserLoginResponse, error)
	ResetPassword(ctx context.Context, user dto.UserResetPassword, forgotPasswordToken string) error
	ForgotPassword(ctx context.Context, user dto.UserForgotPassword) error
	// FetchByParam(ctx context.Context, userParam *dto.UserParam) (dto.UserResponse, error)
	// FetchAll(ctx context.Context, userParam *dto.UserParam, pageParam *dto.PaginationRequest) ([]dto.UserResponse, dto.PaginationResponse, error)
	Logout(ctx context.Context, jwtToken string) error
	DeleteUnverifiedUsers()
}