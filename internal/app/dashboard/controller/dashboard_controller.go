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

	app.Get("/api/v1/admin/dashboard", middleware.Authentication, middleware.RateLimiter(), middleware.AuthAdmin, c.GetAdminDashboard)
	app.Get("/api/v1/organizer/dashboard", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.GetOrganizerDashboard)
	app.Get("/api/v1/peserta/dashboard", middleware.Authentication, middleware.RateLimiter(), middleware.AuthPeserta, c.GetParticipantDashboard)
	app.Get("/api/v1/public/home-stats", middleware.RateLimiter(), c.GetHomeStats)
}

// GetAdminDashboard godoc
// @Summary Get admin dashboard
// @Description Get aggregate statistics, recent top-up transactions, and recent EO registrations for admin dashboard
// @Tags Dashboard
// @Security ApiKeyAuth && BearerAuth
// @Produce json
// @Success 200 {object} response.Response{data=dto.AdminDashboardRes}
// @Router /api/v1/admin/dashboard [get]
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
func (c *dashboardController) GetAdminDashboard(ctx *fiber.Ctx) error {
	res, err := c.svc.GetAdminDashboard(ctx.Context())
	if err != nil {
		response.Send(ctx, http.StatusInternalServerError, "failed to get dashboard", nil, err)
		return nil
	}

	response.Send(ctx, http.StatusOK, "success to get dashboard", res, nil)
	return nil
}

// GetOrganizerDashboard godoc
// @Summary Get organizer dashboard
// @Description Get dashboard statistics and activities for organizer
// @Tags Dashboard
// @Security ApiKeyAuth && BearerAuth
// @Produce json
// @Success 200 {object} response.Response{data=dto.OrganizerDashboardRes}
// @Router /api/v1/organizer/dashboard [get]
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
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
// @Success 200 {object} response.Response{data=dto.ParticipantDashboardRes}
// @Router /api/v1/peserta/dashboard [get]
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
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

// GetHomeStats godoc
// @Summary Get public home stats
// @Description Get aggregate public home page statistics. Events exclude ARCHIVED status. Organizers and participants include active verified accounts. Teams include teams owned by active verified participant accounts.
// @Tags Dashboard
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} response.Response{data=dto.HomeStatsResponse}
// @Router /api/v1/public/home-stats [get]
// @Failure 401 {object} response.ErrorResponse
// @Failure 429 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
func (c *dashboardController) GetHomeStats(ctx *fiber.Ctx) error {
	res, err := c.svc.GetHomeStats(ctx.Context())
	if err != nil {
		response.Send(ctx, http.StatusInternalServerError, "failed to get home stats", nil, err)
		return nil
	}

	response.Send(ctx, http.StatusOK, "success to get home stats", res, nil)
	return nil
}
