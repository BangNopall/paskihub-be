package controller

import (
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/pkg/helpers/http/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ParticipantProfileController struct {
	service  contracts.ParticipantProfileService
	validate *validator.Validate
}

func NewParticipantProfileController(service contracts.ParticipantProfileService, validate *validator.Validate) *ParticipantProfileController {
	return &ParticipantProfileController{
		service:  service,
		validate: validate,
	}
}

func (c *ParticipantProfileController) Route(router fiber.Router) {
	profileGroup := router.Group("/profile")

	profileGroup.Get("", c.GetProfile)
	profileGroup.Put("", c.UpdateInstitution)
	profileGroup.Put("/security", c.UpdatePassword)
}

// GetProfile godoc
// @Summary Get participant profile
// @Description Retrieve current participant's profile
// @Tags Participant Profile
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/peserta/profile [get]
func (c *ParticipantProfileController) GetProfile(ctx *fiber.Ctx) error {
	userID := ctx.Locals("id").(string)

	res, err := c.service.GetProfile(ctx.Context(), userID)
	if err != nil {
		response.Send(ctx, fiber.StatusInternalServerError, "Failed to get profile", nil, err)
		return nil
	}

	response.Send(ctx, fiber.StatusOK, "Successfully retrieved profile", res, nil)
	return nil
}

// UpdateInstitution godoc
// @Summary Update participant institution
// @Description Update the institution of the current participant
// @Tags Participant Profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param req body dto.UpdateInstitutionRequest true "Institution Details"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/peserta/profile [put]
func (c *ParticipantProfileController) UpdateInstitution(ctx *fiber.Ctx) error {
	userID := ctx.Locals("id").(string)

	var req dto.UpdateInstitutionRequest
	if err := ctx.BodyParser(&req); err != nil {
		response.Send(ctx, fiber.StatusBadRequest, "Invalid request body", nil, err)
		return nil
	}

	if err := c.validate.Struct(req); err != nil {
		response.Send(ctx, fiber.StatusBadRequest, "Validation error", nil, err)
		return nil
	}

	err := c.service.UpdateInstitution(ctx.Context(), userID, req)
	if err != nil {
		response.Send(ctx, fiber.StatusInternalServerError, "Failed to update profile", nil, err)
		return nil
	}

	response.Send(ctx, fiber.StatusOK, "Successfully updated profile", nil, nil)
	return nil
}

// UpdatePassword godoc
// @Summary Update participant password
// @Description Update the password for the current participant
// @Tags Participant Profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param req body dto.UpdatePasswordRequest true "Password Details"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/peserta/profile/security [put]
func (c *ParticipantProfileController) UpdatePassword(ctx *fiber.Ctx) error {
	userID := ctx.Locals("id").(string)

	var req dto.UpdatePasswordRequest
	if err := ctx.BodyParser(&req); err != nil {
		response.Send(ctx, fiber.StatusBadRequest, "Invalid request body", nil, err)
		return nil
	}

	if err := c.validate.Struct(req); err != nil {
		response.Send(ctx, fiber.StatusBadRequest, "Validation error", nil, err)
		return nil
	}

	err := c.service.UpdatePassword(ctx.Context(), userID, req)
	if err != nil {
		response.Send(ctx, fiber.StatusInternalServerError, "Failed to update password", nil, err)
		return nil
	}

	response.Send(ctx, fiber.StatusOK, "Successfully updated password", nil, nil)
	return nil
}
