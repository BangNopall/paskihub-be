package controller

import (
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
}

// GetOpenEvents godoc
// @Summary Get open events
// @Description Retrieve a list of events open for registration
// @Tags Participant Event
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/peserta/events/open [get]
func (c *ParticipantEventController) GetOpenEvents(ctx *fiber.Ctx) error {
	res, err := c.service.GetOpenEvents(ctx.Context())
	if err != nil {
		response.Send(ctx, fiber.StatusInternalServerError, "Failed to get open events", nil, err)
		return nil
	}

	response.Send(ctx, fiber.StatusOK, "Successfully retrieved open events", res, nil)
	return nil
}

// RegisterEvent godoc
// @Summary Register to an event
// @Description Register participant team to an event level
// @Tags Participant Event
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param event_level_id formData string true "Event Level ID"
// @Param team_id formData string true "Team ID"
// @Param payment_type formData string true "Payment Type (dp, lunas)"
// @Param payment_proof formData file true "Payment Proof Image"
// @Success 201 {object} map[string]interface{}
// @Router /api/v1/peserta/events/register [post]
func (c *ParticipantEventController) RegisterEvent(ctx *fiber.Ctx) error {
	userID := ctx.Locals("id").(string) // may be used internally for tracing/validation later
	_ = userID 

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

	err = c.service.RegisterEvent(ctx.Context(), req)
	if err != nil {
		response.Send(ctx, fiber.StatusInternalServerError, "Failed to register event", nil, err)
		return nil
	}

	response.Send(ctx, fiber.StatusCreated, "Successfully registered to event", nil, nil)
	return nil
}

// PelunasanEvent godoc
// @Summary Pelunasan payment for event
// @Description Upload full payment proof for an event registration
// @Tags Participant Event
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Registration ID"
// @Param payment_proof formData file true "Payment Proof Image"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/peserta/events/register/{id}/pelunasan [put]
func (c *ParticipantEventController) PelunasanEvent(ctx *fiber.Ctx) error {
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

	err = c.service.PelunasanEvent(ctx.Context(), regisID, req)
	if err != nil {
		response.Send(ctx, fiber.StatusInternalServerError, "Failed to upload pelunasan", nil, err)
		return nil
	}

	response.Send(ctx, fiber.StatusOK, "Successfully uploaded pelunasan", nil, nil)
	return nil
}

// GetActiveEvents godoc
// @Summary Get active events
// @Description Retrieve a list of active registered events for the participant
// @Tags Participant Event
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/peserta/events/active [get]
func (c *ParticipantEventController) GetActiveEvents(ctx *fiber.Ctx) error {
	userID := ctx.Locals("id").(string)

	res, err := c.service.GetActiveEvents(ctx.Context(), userID)
	if err != nil {
		response.Send(ctx, fiber.StatusInternalServerError, "Failed to get active events", nil, err)
		return nil
	}

	response.Send(ctx, fiber.StatusOK, "Successfully retrieved active events", res, nil)
	return nil
}
