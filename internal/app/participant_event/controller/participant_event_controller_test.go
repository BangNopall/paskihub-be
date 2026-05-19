package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type stubParticipantEventService struct{}

func (s stubParticipantEventService) GetOpenEvents(ctx context.Context) ([]dto.OpenEventResponse, error) {
	return nil, nil
}

func (s stubParticipantEventService) RegisterEvent(ctx context.Context, userID string, req dto.RegisterEventRequest) error {
	return nil
}

func (s stubParticipantEventService) PelunasanEvent(ctx context.Context, regisID string, req dto.PelunasanEventRequest) error {
	return nil
}

func (s stubParticipantEventService) GetActiveEvents(ctx context.Context, userID string) ([]dto.ActiveEventResponse, error) {
	return nil, nil
}

func (s stubParticipantEventService) GetRegistrationDetail(ctx context.Context, userID string, regisID string) (dto.RegistrationDetailResponse, error) {
	return dto.RegistrationDetailResponse{EventLevelId: uuid.NewString()}, nil
}

func (s stubParticipantEventService) GetScoreboard(ctx context.Context, eventLevelID string) (dto.ScoreboardResponse, error) {
	return dto.ScoreboardResponse{EventLevelID: uuid.MustParse(eventLevelID)}, nil
}

func TestParticipantEventRouteMountsOverviewEndpoints(t *testing.T) {
	app := fiber.New()
	app.Use(func(ctx *fiber.Ctx) error {
		ctx.Locals("id", uuid.NewString())
		return ctx.Next()
	})

	NewParticipantEventController(stubParticipantEventService{}, nil).Route(app.Group("/api/v1/peserta"))

	regisID := uuid.NewString()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/peserta/events/registrations/"+regisID, nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("request registration detail: %v", err)
	}
	if res.StatusCode == fiber.StatusNotFound {
		t.Fatalf("registration detail route is not mounted")
	}

	eventLevelID := uuid.NewString()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/peserta/rekap/scoreboard/"+eventLevelID, nil)
	res, err = app.Test(req)
	if err != nil {
		t.Fatalf("request scoreboard: %v", err)
	}
	if res.StatusCode == fiber.StatusNotFound {
		t.Fatalf("scoreboard route is not mounted")
	}
}
