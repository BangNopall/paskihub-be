package controller

import (
	"net/http"

	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/internal/middlewares"
	"github.com/BangNopall/paskihub-be/pkg/helpers/http/response"
	"github.com/gofiber/fiber/v2"
)

type juryController struct {
	svc contracts.JuryService
}

func InitJuryController(
	svc contracts.JuryService,
	router fiber.Router,
	middleware *middlewares.Middleware,
) {
	c := &juryController{svc: svc}

	juryRouter := router.Group("/api/v1/organizer/juries")
	juryRouter.Use(middleware.Authentication, middleware.AuthOrganizer)

	juryRouter.Get("/", c.GetAll)
	juryRouter.Post("/", c.Create)
	juryRouter.Put("/:id", c.Update)
	juryRouter.Delete("/:id", c.Delete)
}

// GetAll godoc
// @Summary Get all juries
// @Description Get all juries for organizer
// @Tags Juries
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/organizer/juries [get]
func (c *juryController) GetAll(ctx *fiber.Ctx) error {
	res, err := c.svc.GetAll(ctx.Context())
	if err != nil {
		response.Send(ctx, http.StatusInternalServerError, "failed to get juries", nil, err)
		return nil
	}
	response.Send(ctx, http.StatusOK, "success to get juries", res, nil)
	return nil
}

// Create godoc
// @Summary Create a jury
// @Description Create a new jury
// @Tags Juries
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param jury body dto.JuryRequest true "Jury Request"
// @Success 201 {object} map[string]interface{}
// @Router /api/v1/organizer/juries [post]
func (c *juryController) Create(ctx *fiber.Ctx) error {
	var req dto.JuryRequest
	if err := ctx.BodyParser(&req); err != nil {
		response.Send(ctx, http.StatusBadRequest, "invalid request body", nil, err)
		return nil
	}

	res, err := c.svc.Create(ctx.Context(), req)
	if err != nil {
		response.Send(ctx, http.StatusInternalServerError, "failed to create jury", nil, err)
		return nil
	}
	response.Send(ctx, http.StatusCreated, "success to create jury", res, nil)
	return nil
}

// Update godoc
// @Summary Update a jury
// @Description Update an existing jury
// @Tags Juries
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Jury ID"
// @Param jury body dto.JuryRequest true "Jury Request"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/organizer/juries/{id} [put]
func (c *juryController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req dto.JuryRequest
	if err := ctx.BodyParser(&req); err != nil {
		response.Send(ctx, http.StatusBadRequest, "invalid request body", nil, err)
		return nil
	}

	res, err := c.svc.Update(ctx.Context(), id, req)
	if err != nil {
		response.Send(ctx, http.StatusInternalServerError, "failed to update jury", nil, err)
		return nil
	}
	response.Send(ctx, http.StatusOK, "success to update jury", res, nil)
	return nil
}

// Delete godoc
// @Summary Delete a jury
// @Description Delete a jury by ID
// @Tags Juries
// @Security BearerAuth
// @Produce json
// @Param id path string true "Jury ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/organizer/juries/{id} [delete]
func (c *juryController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if err := c.svc.Delete(ctx.Context(), id); err != nil {
		response.Send(ctx, http.StatusInternalServerError, "failed to delete jury", nil, err)
		return nil
	}
	response.Send(ctx, http.StatusOK, "success to delete jury", nil, nil)
	return nil
}
