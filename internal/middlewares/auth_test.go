package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type authUserRepository struct {
	user entity.User
}

func (r authUserRepository) CreateUser(user *entity.User) error { return nil }
func (r authUserRepository) FindUser(user *entity.User, userParam *dto.UserParam, relations ...string) error {
	*user = r.user
	return nil
}
func (r authUserRepository) UpdateUser(updateUser *dto.UserUpdate, userID uuid.UUID) error {
	return nil
}
func (r authUserRepository) DeleteUnverifiedUser() error { return nil }
func (r authUserRepository) FetchAllUsers(ctx context.Context, role *string) ([]entity.User, error) {
	return nil, nil
}
func (r authUserRepository) UpdateUserStatus(ctx context.Context, userID uuid.UUID, isBanned bool) error {
	return nil
}
func (r authUserRepository) VerifyUserEmail(ctx context.Context, userID uuid.UUID) error {
	return nil
}
func (r authUserRepository) FetchUserDetailAdmin(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	return nil, nil
}
func (r authUserRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error { return nil }

type authJWT struct {
	userID uuid.UUID
}

func (j authJWT) GenerateToken(userID uuid.UUID, payload entity.User) (string, error) {
	return "", nil
}
func (j authJWT) ValidateToken(tokenString string) (uuid.UUID, string, string, *uuid.UUID, error) {
	return j.userID, "banned@example.com", string(enums.Peserta), nil, nil
}

type authRedis struct{}

func (authRedis) Set(ctx context.Context, key string, value interface{}, expDuration time.Duration) error {
	return nil
}
func (authRedis) Get(ctx context.Context, key string) (string, error) { return "", nil }
func (authRedis) Delete(ctx context.Context, key string) error        { return nil }

func TestAuthenticationRejectsTokenForBannedUser(t *testing.T) {
	userID := uuid.New()
	middleware := NewMiddleware(
		authJWT{userID: userID},
		authUserRepository{user: entity.User{Id: userID, IsBanned: true}},
		authRedis{},
	)
	app := fiber.New()
	nextCalled := false
	app.Get("/protected", middleware.Authentication, func(ctx *fiber.Ctx) error {
		nextCalled = true
		return ctx.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer existing-token")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected 403, got %d", res.StatusCode)
	}
	if nextCalled {
		t.Fatal("banned user reached protected handler")
	}
}
