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

type systemSettingController struct {
	settingSvc contracts.SystemSettingService
}

func InitSystemSettingController(
	settingSvc contracts.SystemSettingService,
	router fiber.Router,
	middleware *middlewares.Middleware,
) {
	controller := &systemSettingController{
		settingSvc: settingSvc,
	}

	settingRouter := router.Group("/api/v1/settings")

	// Public Route
	settingRouter.Get("/public", middleware.Authentication, middleware.AuthOrganizer, controller.GetPublicSettings)

	// Admin Routes
	settingRouter.Get("/", middleware.Authentication, middleware.AuthAdmin, controller.GetSettings)
	settingRouter.Patch("/", middleware.Authentication, middleware.AuthAdmin, controller.UpdateSettings)
}

// GetPublicSettings godoc
// @Summary Get public system settings
// @Description Get public system configurations like coin rate and bank info
// @Tags System Settings
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/settings/public [get]
func (c *systemSettingController) GetPublicSettings(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to get public settings"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	res, err = c.settingSvc.GetPublicSettings(ctx.Context())
	code = domain.GetCode(err)
	if err != nil {
		return nil
	}

	message = "success to get public settings"
	return nil
}

// GetSettings godoc
// @Summary Get all system settings (Admin)
// @Description Get all system configurations
// @Tags System Settings
// @Security ApiKeyAuth && BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/settings [get]
func (c *systemSettingController) GetSettings(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to get settings"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	res, err = c.settingSvc.GetSettings(ctx.Context())
	code = domain.GetCode(err)
	if err != nil {
		return nil
	}

	message = "success to get settings"
	return nil
}

// UpdateSettings godoc
// @Summary Update system settings (Admin)
// @Description Update system configurations like coin rate and bank info
// @Tags System Settings
// @Security ApiKeyAuth && BearerAuth
// @Accept json
// @Produce json
// @Param req body dto.UpdateSystemSettingRequest true "Update Settings Request"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/settings [patch]
func (c *systemSettingController) UpdateSettings(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to update settings"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	req := new(dto.UpdateSystemSettingRequest)
	if err := ctx.BodyParser(req); err != nil {
		code = http.StatusBadRequest
		return nil
	}

	err = c.settingSvc.UpdateSettings(ctx.Context(), req)
	code = domain.GetCode(err)
	if err != nil {
		return nil
	}

	message = "success to update settings"
	return nil
}
