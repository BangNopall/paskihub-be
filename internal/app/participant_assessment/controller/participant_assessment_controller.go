package controller

import (
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/pkg/helpers/http/response"
	"github.com/gofiber/fiber/v2"
)

type ParticipantAssessmentController struct {
	service contracts.ParticipantAssessmentService
}

func NewParticipantAssessmentController(service contracts.ParticipantAssessmentService) *ParticipantAssessmentController {
	return &ParticipantAssessmentController{
		service: service,
	}
}

func (c *ParticipantAssessmentController) Route(router fiber.Router) {
	group := router.Group("/assessment")
	group.Get("/recap/:regis_id", c.GetAssessmentRecap)
}

// GetAssessmentRecap godoc
// @Summary Get assessment recap
// @Description Get assessment recap details for a participant registration
// @Tags Participant Assessment
// @Security BearerAuth
// @Produce json
// @Param regis_id path string true "Registration ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/peserta/assessment/recap/{regis_id} [get]
func (c *ParticipantAssessmentController) GetAssessmentRecap(ctx *fiber.Ctx) error {
	userID := ctx.Locals("id").(string)
	regisID := ctx.Params("regis_id")

	res, err := c.service.GetAssessmentRecap(ctx.Context(), userID, regisID)
	if err != nil {
		response.Send(ctx, fiber.StatusInternalServerError, "Failed to get assessment recap", nil, err)
		return nil
	}

	response.Send(ctx, fiber.StatusOK, "Successfully retrieved assessment recap", res, nil)
	return nil
}
