package controller

import (
	"net/http"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/internal/middlewares"
	"github.com/BangNopall/paskihub-be/pkg/helpers/http/response"
	"github.com/gofiber/fiber/v2"
)

type formPenilaianController struct {
	service contracts.FormPenilaianService
}

func InitFormPenilaianController(
	service contracts.FormPenilaianService,
	router fiber.Router,
	middleware *middlewares.Middleware,
) {
	c := &formPenilaianController{service: service}
	routes := router.Group("/api/v1/assessment")

	routes.Post("/scores/bulk", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.BulkInsertScores)
	routes.Post("/violations/bulk", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.BulkInsertViolations)
}

// BulkInsertScores godoc
// @Summary Bulk insert scores
// @Description Insert multiple scores for a team
// @Tags Assessment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param req body dto.BulkInsertScoresRequest true "Bulk Scores Request"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/assessment/scores/bulk [post]
func (c *formPenilaianController) BulkInsertScores(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to insert scores"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	var req dto.BulkInsertScoresRequest
	if err = ctx.BodyParser(&req); err != nil {
		return nil
	}

	err = c.service.BulkInsertScores(ctx.Context(), req)
	code = domain.GetCode(err)
	if err != nil {
		return nil
	}

	message = "success to insert scores"
	return nil
}

// BulkInsertViolations godoc
// @Summary Bulk insert violations
// @Description Insert multiple violations for a team
// @Tags Assessment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param req body dto.BulkInsertViolationsRequest true "Bulk Violations Request"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/assessment/violations/bulk [post]
func (c *formPenilaianController) BulkInsertViolations(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to insert violations"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	var req dto.BulkInsertViolationsRequest
	if err = ctx.BodyParser(&req); err != nil {
		return nil
	}

	err = c.service.BulkInsertTeamViolations(ctx.Context(), req)
	code = domain.GetCode(err)
	if err != nil {
		return nil
	}

	message = "success to insert violations"
	return nil
}
