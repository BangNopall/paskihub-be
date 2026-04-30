package controller

import (
	"net/http"

	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/internal/middlewares"
	"github.com/BangNopall/paskihub-be/pkg/helpers/http/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type dashboardController struct {
	svc contracts.IDashboardService
}

func InitDashboardController(
	svc contracts.IDashboardService,
	app *fiber.App,
	middleware *middlewares.Middleware,
) {
	c := &dashboardController{svc: svc}

	app.Get("/api/v1/organizer/dashboard", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.GetOrganizerDashboard)
	app.Get("/api/v1/peserta/dashboard", middleware.Authentication, middleware.RateLimiter(), middleware.AuthPeserta, c.GetParticipantDashboard)
}

// GetOrganizerDashboard godoc
// @Summary Get organizer dashboard
// @Description Get dashboard statistics and activities for organizer
// @Tags Dashboard
// @Security ApiKeyAuth && BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/organizer/dashboard [get]
func (c *dashboardController) GetOrganizerDashboard(ctx *fiber.Ctx) error {
	userIdStr := ctx.Locals("id").(string)
	userId, _ := uuid.Parse(userIdStr)

	res, err := c.svc.GetOrganizerDashboard(ctx.Context(), userId)
	if err != nil {
		response.Send(ctx, http.StatusInternalServerError, "failed to get dashboard", nil, err)
		return nil
	}

	response.Send(ctx, http.StatusOK, "success to get dashboard", res, nil)
	return nil
}

// GetParticipantDashboard godoc
// @Summary Get participant dashboard
// @Description Get dashboard statistics and activities for participant
// @Tags Dashboard
// @Security ApiKeyAuth && BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/peserta/dashboard [get]
func (c *dashboardController) GetParticipantDashboard(ctx *fiber.Ctx) error {
	userIdStr := ctx.Locals("id").(string)
	userId, _ := uuid.Parse(userIdStr)

	res, err := c.svc.GetParticipantDashboard(ctx.Context(), userId)
	if err != nil {
		response.Send(ctx, http.StatusInternalServerError, "failed to get dashboard", nil, err)
		return nil
	}

	response.Send(ctx, http.StatusOK, "success to get dashboard", res, nil)
	return nil
}
