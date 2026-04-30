package controller

import (
	"net/http"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/internal/middlewares"
	"github.com/BangNopall/paskihub-be/pkg/helpers/http/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type eoController struct {
	svc       contracts.IEOService
	validator *validator.Validate
}

func NewEOController(svc contracts.IEOService, v *validator.Validate) *eoController {
	return &eoController{svc: svc, validator: v}
}

func (c *eoController) Route(router fiber.Router, middleware *middlewares.Middleware) {
	group := router.Group("/api/v1/eo", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer)

	group.Get("/profile", c.GetProfile)
	group.Put("/profile/password", c.UpdatePassword)
	group.Get("/staff", c.GetStaffs)
	group.Post("/staff", c.CreateStaff)
	group.Put("/staff/:id/password", c.ResetStaffPassword)
	group.Delete("/staff/:id", c.DeleteStaff)
}

func getUserId(ctx *fiber.Ctx) (uuid.UUID, error) {
	userIdStr := ctx.Locals("id").(string)
	return uuid.Parse(userIdStr)
}

// GetProfile godoc
// @Summary Get EO profile
// @Description Retrieve current Event Organizer's profile information
// @Tags EO Profile
// @Security ApiKeyAuth && BearerAuth
// @Produce json
// @Success 200 {object} response.Response{data=dto.EOProfileRes}
// @Router /api/v1/eo/profile [get]
func (c *eoController) GetProfile(ctx *fiber.Ctx) error {
	userId, err := getUserId(ctx)
	if err != nil {
		response.SendErrResp(ctx, http.StatusUnauthorized, response.Fail, "unauthorized", err)
		return nil
	}

	res, err := c.svc.GetProfile(ctx.Context(), userId)
	if err != nil {
		response.SendErrResp(ctx, domain.GetCode(err), response.Fail, "failed to get profile", err)
		return nil
	}

	response.Send(ctx, http.StatusOK, "success to get profile", res, nil)
	return nil
}

// UpdatePassword godoc
// @Summary Update EO password
// @Description Update password for the current Event Organizer account
// @Tags EO Profile
// @Security ApiKeyAuth && BearerAuth
// @Accept json
// @Produce json
// @Param req body dto.EOUpdatePasswordReq true "Update Password Request"
// @Success 200 {object} response.Response
// @Router /api/v1/eo/profile/password [put]
func (c *eoController) UpdatePassword(ctx *fiber.Ctx) error {
	userId, err := getUserId(ctx)
	if err != nil {
		response.SendErrResp(ctx, http.StatusUnauthorized, response.Fail, "unauthorized", err)
		return nil
	}

	var req dto.EOUpdatePasswordReq
	if err := ctx.BodyParser(&req); err != nil {
		response.SendErrResp(ctx, http.StatusBadRequest, response.Fail, "invalid request body", err)
		return nil
	}

	if err := c.validator.Struct(req); err != nil {
		response.SendErrResp(ctx, http.StatusBadRequest, response.Fail, "validation failed", err)
		return nil
	}

	err = c.svc.UpdatePassword(ctx.Context(), userId, req)
	if err != nil {
		response.SendErrResp(ctx, domain.GetCode(err), response.Fail, "failed to update password", err)
		return nil
	}

	response.Send(ctx, http.StatusOK, "success to update password", nil, nil)
	return nil
}

// GetStaffs godoc
// @Summary Get list of staff
// @Description Retrieve a list of staff accounts associated with the current Event Organizer
// @Tags EO Staff
// @Security ApiKeyAuth && BearerAuth
// @Produce json
// @Success 200 {object} response.Response{data=[]dto.EOStaffRes}
// @Router /api/v1/eo/staff [get]
func (c *eoController) GetStaffs(ctx *fiber.Ctx) error {
	userId, err := getUserId(ctx)
	if err != nil {
		response.SendErrResp(ctx, http.StatusUnauthorized, response.Fail, "unauthorized", err)
		return nil
	}

	res, err := c.svc.GetStaffs(ctx.Context(), userId)
	if err != nil {
		response.SendErrResp(ctx, domain.GetCode(err), response.Fail, "failed to get staffs", err)
		return nil
	}

	response.Send(ctx, http.StatusOK, "success to get staffs", res, nil)
	return nil
}

// CreateStaff godoc
// @Summary Create a new staff account
// @Description Creates a new staff account under the current Event Organizer
// @Tags EO Staff
// @Security ApiKeyAuth && BearerAuth
// @Accept json
// @Produce json
// @Param req body dto.EOStaffCreateReq true "Staff Create Request"
// @Success 201 {object} response.Response
// @Router /api/v1/eo/staff [post]
func (c *eoController) CreateStaff(ctx *fiber.Ctx) error {
	userId, err := getUserId(ctx)
	if err != nil {
		response.SendErrResp(ctx, http.StatusUnauthorized, response.Fail, "unauthorized", err)
		return nil
	}

	var req dto.EOStaffCreateReq
	if err := ctx.BodyParser(&req); err != nil {
		response.SendErrResp(ctx, http.StatusBadRequest, response.Fail, "invalid request body", err)
		return nil
	}

	if err := c.validator.Struct(req); err != nil {
		response.SendErrResp(ctx, http.StatusBadRequest, response.Fail, "validation failed", err)
		return nil
	}

	err = c.svc.CreateStaff(ctx.Context(), userId, req)
	if err != nil {
		response.SendErrResp(ctx, domain.GetCode(err), response.Fail, "failed to create staff", err)
		return nil
	}

	response.Send(ctx, http.StatusCreated, "success to create staff", nil, nil)
	return nil
}

// ResetStaffPassword godoc
// @Summary Reset staff password
// @Description Updates the password for a specific staff account (Reset by EO)
// @Tags EO Staff
// @Security ApiKeyAuth && BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Staff ID"
// @Param req body dto.EOStaffResetPasswordReq true "Reset Password Request"
// @Success 200 {object} response.Response
// @Router /api/v1/eo/staff/{id}/password [put]
func (c *eoController) ResetStaffPassword(ctx *fiber.Ctx) error {
	parentId, err := getUserId(ctx)
	if err != nil {
		response.SendErrResp(ctx, http.StatusUnauthorized, response.Fail, "unauthorized", err)
		return nil
	}

	staffIdStr := ctx.Params("id")
	staffId, err := uuid.Parse(staffIdStr)
	if err != nil {
		response.SendErrResp(ctx, http.StatusBadRequest, response.Fail, "invalid staff id", err)
		return nil
	}

	var req dto.EOStaffResetPasswordReq
	if err := ctx.BodyParser(&req); err != nil {
		response.SendErrResp(ctx, http.StatusBadRequest, response.Fail, "invalid request body", err)
		return nil
	}

	if err := c.validator.Struct(req); err != nil {
		response.SendErrResp(ctx, http.StatusBadRequest, response.Fail, "validation failed", err)
		return nil
	}

	err = c.svc.ResetStaffPassword(ctx.Context(), parentId, staffId, req)
	if err != nil {
		response.SendErrResp(ctx, domain.GetCode(err), response.Fail, "failed to reset staff password", err)
		return nil
	}

	response.Send(ctx, http.StatusOK, "success to reset staff password", nil, nil)
	return nil
}

// DeleteStaff godoc
// @Summary Delete a staff account
// @Description Permanently removes a staff account
// @Tags EO Staff
// @Security ApiKeyAuth && BearerAuth
// @Produce json
// @Param id path string true "Staff ID"
// @Success 200 {object} response.Response
// @Router /api/v1/eo/staff/{id} [delete]
func (c *eoController) DeleteStaff(ctx *fiber.Ctx) error {
	parentId, err := getUserId(ctx)
	if err != nil {
		response.SendErrResp(ctx, http.StatusUnauthorized, response.Fail, "unauthorized", err)
		return nil
	}

	staffIdStr := ctx.Params("id")
	staffId, err := uuid.Parse(staffIdStr)
	if err != nil {
		response.SendErrResp(ctx, http.StatusBadRequest, response.Fail, "invalid staff id", err)
		return nil
	}

	err = c.svc.DeleteStaff(ctx.Context(), parentId, staffId)
	if err != nil {
		response.SendErrResp(ctx, domain.GetCode(err), response.Fail, "failed to delete staff", err)
		return nil
	}

	response.Send(ctx, http.StatusOK, "success to delete staff", nil, nil)
	return nil
}
