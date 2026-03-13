package controller

import (
	"net/http"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/internal/middlewares"
	"github.com/BangNopall/paskihub-be/pkg/helpers/http/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type rekapController struct {
	service contracts.RekapService
}

func InitRekapController(
	service contracts.RekapService,
	router fiber.Router,
	middleware *middlewares.Middleware,
) {
	c := &rekapController{service: service}
	routes := router.Group("/api/v1/rekap")

	routes.Get("/detail/:regisId", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.GetTeamAssessmentDetail)
	routes.Get("/scoreboard/:eventLevelId", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.GetScoreboard)
	routes.Post("/leaderboard/custom/:eventLevelId", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.GetLeaderboardCustom)
	routes.Put("/publish/:eventLevelId", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.PublishScoreboard)
}

func (c *rekapController) GetTeamAssessmentDetail(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to get detail"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	regisIdStr := ctx.Params("regisId")
	regisId, err := uuid.Parse(regisIdStr)
	if err != nil {
		return nil
	}

	res, err = c.service.GetTeamAssessmentDetail(ctx.Context(), regisId)
	code = domain.GetCode(err)
	if err != nil {
		return nil
	}

	message = "success to get detail"
	return nil
}

func (c *rekapController) GetScoreboard(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to get scoreboard"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	eventLevelIdStr := ctx.Params("eventLevelId")
	eventLevelId, err := uuid.Parse(eventLevelIdStr)
	if err != nil {
		return nil
	}

	res, err = c.service.GetScoreboard(ctx.Context(), eventLevelId)
	code = domain.GetCode(err)
	if err != nil {
		return nil
	}

	message = "success to get scoreboard"
	return nil
}

func (c *rekapController) GetLeaderboardCustom(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to get leaderboard"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	eventLevelIdStr := ctx.Params("eventLevelId")
	eventLevelId, err := uuid.Parse(eventLevelIdStr)
	if err != nil {
		return nil
	}

	var req dto.CustomLeaderboardRequest
	if err = ctx.BodyParser(&req); err != nil {
		return nil
	}

	res, err = c.service.GetLeaderboardCustom(ctx.Context(), req, eventLevelId)
	code = domain.GetCode(err)
	if err != nil {
		return nil
	}

	message = "success to get leaderboard"
	return nil
}

func (c *rekapController) PublishScoreboard(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to publish scoreboard"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	eventLevelIdStr := ctx.Params("eventLevelId")
	eventLevelId, err := uuid.Parse(eventLevelIdStr)
	if err != nil {
		return nil
	}

	userIdStr := ctx.Locals("id").(string)
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		code = http.StatusUnauthorized
		return nil
	}

	var req dto.PublishScoreboardRequest
	if err = ctx.BodyParser(&req); err != nil {
		return nil
	}

	err = c.service.PublishScoreboard(ctx.Context(), req, eventLevelId, userId)
	code = domain.GetCode(err)
	if err != nil {
		return nil
	}

	message = "success to publish scoreboard"
	return nil
}
