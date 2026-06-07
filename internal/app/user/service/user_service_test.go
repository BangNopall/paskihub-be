package service

import (
	"context"
	"testing"
	"time"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/google/uuid"
)

type stubUserRepository struct {
	user entity.User
}

func (r stubUserRepository) CreateUser(user *entity.User) error { return nil }
func (r stubUserRepository) FindUser(user *entity.User, userParam *dto.UserParam, relations ...string) error {
	*user = r.user
	return nil
}
func (r stubUserRepository) UpdateUser(updateUser *dto.UserUpdate, userID uuid.UUID) error {
	return nil
}
func (r stubUserRepository) DeleteUnverifiedUser() error { return nil }
func (r stubUserRepository) FetchAllUsers(ctx context.Context, role *string) ([]entity.User, error) {
	return nil, nil
}
func (r stubUserRepository) UpdateUserStatus(ctx context.Context, userID uuid.UUID, isBanned bool) error {
	return nil
}
func (r stubUserRepository) VerifyUserEmail(ctx context.Context, userID uuid.UUID) error {
	return nil
}
func (r stubUserRepository) FetchUserDetailAdmin(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	return nil, nil
}
func (r stubUserRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return nil
}

type stubBcrypt struct{}

func (stubBcrypt) Hash(plain string) (string, error) { return plain, nil }
func (stubBcrypt) Compare(password, hashed string) bool {
	return true
}

type stubJWT struct {
	generateCalled *bool
}

func (j stubJWT) GenerateToken(userID uuid.UUID, payload entity.User) (string, error) {
	if j.generateCalled != nil {
		*j.generateCalled = true
	}
	return "token", nil
}
func (stubJWT) ValidateToken(tokenString string) (uuid.UUID, string, string, *uuid.UUID, error) {
	return uuid.Nil, "", "", nil, nil
}

type stubTime struct{}

func (stubTime) Now() time.Time                       { return time.Now() }
func (stubTime) Add(duration time.Duration) time.Time { return time.Now().Add(duration) }

func TestLoginRejectsBannedUserBeforeTokenGeneration(t *testing.T) {
	generateCalled := false
	svc := &userService{
		userRepo: stubUserRepository{user: entity.User{
			Id:              uuid.New(),
			Email:           "banned@example.com",
			Password:        "hash",
			Role:            enums.Peserta,
			EmailIsVerified: true,
			IsBanned:        true,
		}},
		bcrypt:  stubBcrypt{},
		jwt:     stubJWT{generateCalled: &generateCalled},
		time:    stubTime{},
		timeout: time.Second,
	}

	_, err := svc.Login(context.Background(), dto.UserLogin{
		Email:    "banned@example.com",
		Password: "password",
	})

	if err != domain.ErrAccountBanned {
		t.Fatalf("expected ErrAccountBanned, got %v", err)
	}
	if generateCalled {
		t.Fatal("token was generated for a banned user")
	}
}
