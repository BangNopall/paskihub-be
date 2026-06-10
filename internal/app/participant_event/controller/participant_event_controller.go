package controller

import (
	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/pkg/helpers/http/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ParticipantEventController struct {
	service  contracts.ParticipantEventService
	validate *validator.Validate
}

func NewParticipantEventController(service contracts.ParticipantEventService, validate *validator.Validate) *ParticipantEventController {
	return &ParticipantEventController{
		service:  service,
		validate: validate,
	}
}

func (c *ParticipantEventController) Route(router fiber.Router) {
	eventGroup := router.Group("/events")

	eventGroup.Get("/open", c.GetOpenEvents)
	eventGroup.Post("/register", c.RegisterEvent)
	eventGroup.Put("/register/:id/pelunasan", c.PelunasanEvent)
	eventGroup.Get("/active", c.GetActiveEvents)
	eventGroup.Get("/registrations/:regis_id", c.GetRegistrationDetail)

	rekapGroup := router.Group("/rekap")
	rekapGroup.Get("/scoreboard/:event_level_id", c.GetScoreboard)
}

// GetOpenEvents godoc
// @Summary Get open events
// @Description Retrieve a list of events open for registration
// @Tags Participant Event
// @Security ApiKeyAuth && BearerAuth
// @Produce json
// @Param location query string false "Filter by location"
// @Param search query string false "Search by event name"
// @Success 200 {object} response.Response{data=[]dto.OpenEventResponse}
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/peserta/events/open [get]
func (c *ParticipantEventController) GetOpenEvents(ctx *fiber.Ctx) error {
	location := ctx.Query("location")
	search := ctx.Query("search")
	res, err := c.service.GetOpenEvents(ctx.Context(), location, search)
	if err != nil {
		response.Send(ctx, domain.GetCode(err), "Failed to get open events", nil, err)
		return nil
	}

	response.Send(ctx, fiber.StatusOK, "Successfully retrieved open events", res, nil)
	return nil
}

// RegisterEvent godoc
// @Summary Register to an event
// @Description Register participant team to an event level
// @Tags Participant Event
// @Security ApiKeyAuth && BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param event_level_id formData string true "Event Level ID"
// @Param team_id formData string true "Team ID"
// @Param payment_type formData string true "Payment Type (dp, lunas)"
// @Param payment_proof formData file true "Payment Proof Image"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/peserta/events/register [post]
func (c *ParticipantEventController) RegisterEvent(ctx *fiber.Ctx) error {
	userID := ctx.Locals("id").(string)

	form, err := ctx.MultipartForm()
	if err != nil {
		response.Send(ctx, fiber.StatusBadRequest, "Failed to parse form", nil, err)
		return nil
	}

	var req dto.RegisterEventRequest
	if lvlVals, ok := form.Value["event_level_id"]; ok && len(lvlVals) > 0 {
		req.EventLevelId = lvlVals[0]
	}
	if teamVals, ok := form.Value["team_id"]; ok && len(teamVals) > 0 {
		req.TeamId = teamVals[0]
	}
	if typeVals, ok := form.Value["payment_type"]; ok && len(typeVals) > 0 {
		req.PaymentType = typeVals[0]
	}

	if proofFiles, ok := form.File["payment_proof"]; ok && len(proofFiles) > 0 {
		req.PaymentProof = proofFiles[0]
	}

	if err := c.validate.Struct(req); err != nil {
		response.Send(ctx, fiber.StatusBadRequest, "Validation error", nil, err)
		return nil
	}

	err = c.service.RegisterEvent(ctx.Context(), userID, req)
	if err != nil {
		response.Send(ctx, domain.GetCode(err), "Failed to register event", nil, err)
		return nil
	}

	response.Send(ctx, fiber.StatusCreated, "Successfully registered to event", nil, nil)
	return nil
}

// PelunasanEvent godoc
// @Summary Pelunasan payment for event
// @Description Upload full payment proof for an event registration
// @Tags Participant Event
// @Security ApiKeyAuth && BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Registration ID"
// @Param payment_proof formData file true "Payment Proof Image"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/peserta/events/register/{id}/pelunasan [put]
func (c *ParticipantEventController) PelunasanEvent(ctx *fiber.Ctx) error {
	userID := ctx.Locals("id").(string)
	regisID := ctx.Params("id")

	form, err := ctx.MultipartForm()
	if err != nil {
		response.Send(ctx, fiber.StatusBadRequest, "Failed to parse form", nil, err)
		return nil
	}

	var req dto.PelunasanEventRequest
	if proofFiles, ok := form.File["payment_proof"]; ok && len(proofFiles) > 0 {
		req.PaymentProof = proofFiles[0]
	}

	err = c.service.PelunasanEvent(ctx.Context(), userID, regisID, req)
	if err != nil {
		response.Send(ctx, domain.GetCode(err), "Failed to upload pelunasan", nil, err)
		return nil
	}

	response.Send(ctx, fiber.StatusOK, "Successfully uploaded pelunasan", nil, nil)
	return nil
}

// GetActiveEvents godoc
// @Summary Get active events
// @Description Retrieve a list of active registered events for the participant
// @Tags Participant Event
// @Security ApiKeyAuth && BearerAuth
// @Produce json
// @Success 200 {object} response.Response{data=[]dto.ActiveEventResponse}
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/peserta/events/active [get]
func (c *ParticipantEventController) GetActiveEvents(ctx *fiber.Ctx) error {
	userID := ctx.Locals("id").(string)

	res, err := c.service.GetActiveEvents(ctx.Context(), userID)
	if err != nil {
		response.Send(ctx, domain.GetCode(err), "Failed to get active events", nil, err)
		return nil
	}

	response.Send(ctx, fiber.StatusOK, "Successfully retrieved active events", res, nil)
	return nil
}

// GetRegistrationDetail godoc
// @Summary Get registration detail
// @Description Get specific registration details including event, team and payment
// @Tags Participant Event
// @Security ApiKeyAuth && BearerAuth
// @Produce json
// @Param regis_id path string true "Registration ID"
// @Success 200 {object} response.Response{data=dto.RegistrationDetailResponse}
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/peserta/events/registrations/{regis_id} [get]
func (c *ParticipantEventController) GetRegistrationDetail(ctx *fiber.Ctx) error {
	userID := ctx.Locals("id").(string)
	regisID := ctx.Params("regis_id")

	res, err := c.service.GetRegistrationDetail(ctx.Context(), userID, regisID)
	if err != nil {
		response.Send(ctx, domain.GetCode(err), "Failed to get registration detail", nil, err)
		return nil
	}

	response.Send(ctx, fiber.StatusOK, "Successfully retrieved registration detail", res, nil)
	return nil
}

// GetScoreboard godoc
// @Summary Get scoreboard for participants
// @Description Retrieve the scoreboard if published by the organizer
// @Tags Participant Event
// @Security ApiKeyAuth && BearerAuth
// @Produce json
// @Param event_level_id path string true "Event Level ID"
// @Success 200 {object} response.Response{data=dto.ScoreboardResponse}
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/peserta/rekap/scoreboard/{event_level_id} [get]
func (c *ParticipantEventController) GetScoreboard(ctx *fiber.Ctx) error {
	eventLevelID := ctx.Params("event_level_id")

	res, err := c.service.GetScoreboard(ctx.Context(), eventLevelID)
	if err != nil {
		response.Send(ctx, domain.GetCode(err), "Failed to get scoreboard", nil, err)
		return nil
	}

	response.Send(ctx, fiber.StatusOK, "Successfully retrieved scoreboard", res, nil)
	return nil
}
