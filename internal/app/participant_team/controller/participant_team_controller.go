package controller

import (
	"fmt"

	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/pkg/helpers/http/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ParticipantTeamController struct {
	service  contracts.ParticipantTeamService
	validate *validator.Validate
}

func NewParticipantTeamController(service contracts.ParticipantTeamService, validate *validator.Validate) *ParticipantTeamController {
	return &ParticipantTeamController{
		service:  service,
		validate: validate,
	}
}

func (c *ParticipantTeamController) Route(router fiber.Router) {
	teamGroup := router.Group("/teams")

	teamGroup.Post("", c.CreateTeam)
	teamGroup.Get("", c.GetTeams)
	teamGroup.Get("/:id", c.GetTeamDetail)
	teamGroup.Delete("/:id", c.DeleteTeam)
}

// CreateTeam godoc
// @Summary Create a new participant team
// @Description Upload team details including logo, recommendation letter, and members
// @Tags Participant Team
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Team Name"
// @Param pelatih_name formData string true "Coach Name"
// @Param logo_team formData file true "Team Logo"
// @Param surat_rekomendasi formData file true "Recommendation Letter"
// @Success 201 {object} map[string]interface{}
// @Router /api/v1/peserta/teams [post]
func (c *ParticipantTeamController) CreateTeam(ctx *fiber.Ctx) error {
	userID := ctx.Locals("id").(string)

	form, err := ctx.MultipartForm()
	if err != nil {
		response.Send(ctx, fiber.StatusBadRequest, "Failed to parse multipart form", nil, err)
		return nil
	}

	var req dto.CreateTeamRequest

	if nameVals, ok := form.Value["name"]; ok && len(nameVals) > 0 {
		req.Name = nameVals[0]
	}
	if pelatihVals, ok := form.Value["pelatih_name"]; ok && len(pelatihVals) > 0 {
		req.PelatihName = pelatihVals[0]
	}

	if logoFiles, ok := form.File["logo_team"]; ok && len(logoFiles) > 0 {
		req.LogoTeam = logoFiles[0]
	}
	if recFiles, ok := form.File["surat_rekomendasi"]; ok && len(recFiles) > 0 {
		req.SuratRekomendasi = recFiles[0]
	}

	for i := 0; i < 50; i++ {
		fullNameKey := fmt.Sprintf("members[%d][full_name]", i)
		roleKey := fmt.Sprintf("members[%d][role]", i)
		idCardKey := fmt.Sprintf("members[%d][id_card]", i)
		photoKey := fmt.Sprintf("members[%d][photo]", i)

		var member dto.TeamMemberRequest
		hasData := false

		if fnVals, ok := form.Value[fullNameKey]; ok && len(fnVals) > 0 {
			member.FullName = fnVals[0]
			hasData = true
		}
		if rlVals, ok := form.Value[roleKey]; ok && len(rlVals) > 0 {
			member.Role = rlVals[0]
			hasData = true
		}
		if icFiles, ok := form.File[idCardKey]; ok && len(icFiles) > 0 {
			member.IdCard = icFiles[0]
			hasData = true
		}
		if phFiles, ok := form.File[photoKey]; ok && len(phFiles) > 0 {
			member.Photo = phFiles[0]
			hasData = true
		}

		if !hasData {
			break
		}
		req.Members = append(req.Members, member)
	}

	if err := c.validate.Struct(req); err != nil {
		response.Send(ctx, fiber.StatusBadRequest, "Validation error", nil, err)
		return nil
	}

	err = c.service.CreateTeam(ctx.Context(), userID, req)
	if err != nil {
		response.Send(ctx, fiber.StatusInternalServerError, "Failed to create team", nil, err)
		return nil
	}

	response.Send(ctx, fiber.StatusCreated, "Successfully created team", nil, nil)
	return nil
}

// GetTeams godoc
// @Summary Get participant teams
// @Description Retrieve a list of teams for the current participant
// @Tags Participant Team
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/peserta/teams [get]
func (c *ParticipantTeamController) GetTeams(ctx *fiber.Ctx) error {
	userID := ctx.Locals("id").(string)

	res, err := c.service.GetTeams(ctx.Context(), userID)
	if err != nil {
		response.Send(ctx, fiber.StatusInternalServerError, "Failed to get teams", nil, err)
		return nil
	}

	response.Send(ctx, fiber.StatusOK, "Successfully retrieved teams", res, nil)
	return nil
}

// GetTeamDetail godoc
// @Summary Get team details
// @Description Retrieve details of a specific team
// @Tags Participant Team
// @Security BearerAuth
// @Produce json
// @Param id path string true "Team ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/peserta/teams/{id} [get]
func (c *ParticipantTeamController) GetTeamDetail(ctx *fiber.Ctx) error {
	userID := ctx.Locals("id").(string)
	teamID := ctx.Params("id")

	res, err := c.service.GetTeamDetail(ctx.Context(), userID, teamID)
	if err != nil {
		response.Send(ctx, fiber.StatusInternalServerError, "Failed to get team detail", nil, err)
		return nil
	}

	response.Send(ctx, fiber.StatusOK, "Successfully retrieved team detail", res, nil)
	return nil
}

// DeleteTeam godoc
// @Summary Delete participant team
// @Description Remove a team by its ID
// @Tags Participant Team
// @Security BearerAuth
// @Produce json
// @Param id path string true "Team ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/peserta/teams/{id} [delete]
func (c *ParticipantTeamController) DeleteTeam(ctx *fiber.Ctx) error {
	userID := ctx.Locals("id").(string)
	teamID := ctx.Params("id")

	err := c.service.DeleteTeam(ctx.Context(), userID, teamID)
	if err != nil {
		response.Send(ctx, fiber.StatusBadRequest, "Failed to delete team", nil, err)
		return nil
	}

	response.Send(ctx, fiber.StatusOK, "Successfully deleted team", nil, nil)
	return nil
}
