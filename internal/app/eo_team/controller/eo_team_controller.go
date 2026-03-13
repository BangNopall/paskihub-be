package controller

import (
	"net/http"

	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/internal/middlewares"
	"github.com/BangNopall/paskihub-be/pkg/helpers/http/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type eoTeamController struct {
	svc contracts.IEOTeamService
}

func InitEOTeamController(
	svc contracts.IEOTeamService,
	router fiber.Router,
	middleware *middlewares.Middleware,
) {
	c := &eoTeamController{svc: svc}

	group := router.Group("/api/v1/eo/events/:eventId/teams")

	group.Get("/", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.GetListTeam)
	group.Get("/:registrationId", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.GetDetailTeam)
	group.Put("/:registrationId/approve", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.ApproveTeam)
	group.Put("/:registrationId/reject", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.RejectTeam)
	group.Put("/:registrationId/kick", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.KickTeam)
}

func getUUIDParam(ctx *fiber.Ctx, key string) (uuid.UUID, error) {
	val := ctx.Params(key)
	return uuid.Parse(val)
}

func getUserId(ctx *fiber.Ctx) (uuid.UUID, error) {
	userIdStr := ctx.Locals("id").(string)
	return uuid.Parse(userIdStr)
}

func setCodeFromError(err error) int {
	if err == nil {
		return http.StatusOK
	}
	// For simplicity, returning 400 for generic errors, 404 for not found, 401/403 for unauthorized
	if err.Error() == "unauthorized: you do not own this event" {
		return http.StatusForbidden
	}
	if err.Error() == "registration not found for this event" {
		return http.StatusNotFound
	}
	return http.StatusBadRequest
}

// GetListTeam godoc
// @Summary Get list of teams
// @Description Retrieve a list of team registrations for an event (Event Organizer)
// @Tags EO Team
// @Security BearerAuth
// @Produce json
// @Param eventId path string true "Event ID"
// @Param event_level_id query string false "Event Level ID filter"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/eo/events/{eventId}/teams [get]
func (c *eoTeamController) GetListTeam(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusOK
		res     interface{}
		message string = "failed to get list team"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	userId, err := getUserId(ctx)
	if err != nil {
		code = http.StatusUnauthorized
		return nil
	}

	eventId, err := getUUIDParam(ctx, "eventId")
	if err != nil {
		code = http.StatusBadRequest
		return nil
	}

	var eventLevelId *uuid.UUID
	if elid := ctx.Query("event_level_id"); elid != "" {
		if parsedId, errParse := uuid.Parse(elid); errParse == nil {
			eventLevelId = &parsedId
		}
	}

	res, err = c.svc.GetListTeam(ctx.Context(), eventId, userId, eventLevelId)
	code = setCodeFromError(err)
	if err == nil {
		message = "success to get list team"
	}
	return nil
}

// GetDetailTeam godoc
// @Summary Get team details
// @Description Retrieve details of a specific team registration
// @Tags EO Team
// @Security BearerAuth
// @Produce json
// @Param eventId path string true "Event ID"
// @Param registrationId path string true "Registration ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/eo/events/{eventId}/teams/{registrationId} [get]
func (c *eoTeamController) GetDetailTeam(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusOK
		res     interface{}
		message string = "failed to get detail team"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	userId, err := getUserId(ctx)
	if err != nil {
		code = http.StatusUnauthorized
		return nil
	}

	eventId, err := getUUIDParam(ctx, "eventId")
	if err != nil {
		code = http.StatusBadRequest
		return nil
	}

	registrationId, err := getUUIDParam(ctx, "registrationId")
	if err != nil {
		code = http.StatusBadRequest
		return nil
	}

	res, err = c.svc.GetDetailTeam(ctx.Context(), eventId, userId, registrationId)
	code = setCodeFromError(err)
	if err == nil {
		message = "success to get detail team"
	}
	return nil
}

// ApproveTeam godoc
// @Summary Approve team registration
// @Description Approve a team's registration
// @Tags EO Team
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Param registrationId path string true "Registration ID"
// @Param req body dto.EOTeamApproveReq true "Approval Details"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/eo/events/{eventId}/teams/{registrationId}/approve [put]
func (c *eoTeamController) ApproveTeam(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusOK
		res     interface{}
		message string = "failed to approve team"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	userId, err := getUserId(ctx)
	if err != nil {
		code = http.StatusUnauthorized
		return nil
	}

	eventId, err := getUUIDParam(ctx, "eventId")
	if err != nil {
		code = http.StatusBadRequest
		return nil
	}

	registrationId, err := getUUIDParam(ctx, "registrationId")
	if err != nil {
		code = http.StatusBadRequest
		return nil
	}

	var req dto.EOTeamApproveReq
	if errParse := ctx.BodyParser(&req); errParse != nil {
		err = errParse
		code = http.StatusBadRequest
		return nil
	}

	err = c.svc.ApproveTeam(ctx.Context(), eventId, userId, registrationId, req)
	code = setCodeFromError(err)
	if err == nil {
		message = "success to approve team"
	}
	return nil
}

// RejectTeam godoc
// @Summary Reject team registration
// @Description Reject a pending team registration
// @Tags EO Team
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Param registrationId path string true "Registration ID"
// @Param req body dto.EOTeamRejectReq true "Rejection Reason"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/eo/events/{eventId}/teams/{registrationId}/reject [put]
func (c *eoTeamController) RejectTeam(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusOK
		res     interface{}
		message string = "failed to reject team"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	userId, err := getUserId(ctx)
	if err != nil {
		code = http.StatusUnauthorized
		return nil
	}

	eventId, err := getUUIDParam(ctx, "eventId")
	if err != nil {
		code = http.StatusBadRequest
		return nil
	}

	registrationId, err := getUUIDParam(ctx, "registrationId")
	if err != nil {
		code = http.StatusBadRequest
		return nil
	}

	var req dto.EOTeamRejectReq
	if errParse := ctx.BodyParser(&req); errParse != nil {
		err = errParse
		code = http.StatusBadRequest
		return nil
	}

	err = c.svc.RejectTeam(ctx.Context(), eventId, userId, registrationId, req)
	code = setCodeFromError(err)
	if err == nil {
		message = "success to reject team"
	}
	return nil
}

// KickTeam godoc
// @Summary Kick team
// @Description Disqualify or kick an approved team
// @Tags EO Team
// @Security BearerAuth
// @Produce json
// @Param eventId path string true "Event ID"
// @Param registrationId path string true "Registration ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/eo/events/{eventId}/teams/{registrationId}/kick [put]
func (c *eoTeamController) KickTeam(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusOK
		res     interface{}
		message string = "failed to kick team"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	userId, err := getUserId(ctx)
	if err != nil {
		code = http.StatusUnauthorized
		return nil
	}

	eventId, err := getUUIDParam(ctx, "eventId")
	if err != nil {
		code = http.StatusBadRequest
		return nil
	}

	registrationId, err := getUUIDParam(ctx, "registrationId")
	if err != nil {
		code = http.StatusBadRequest
		return nil
	}

	err = c.svc.KickTeam(ctx.Context(), eventId, userId, registrationId)
	code = setCodeFromError(err)
	if err == nil {
		message = "success to kick team"
	}
	return nil
}
