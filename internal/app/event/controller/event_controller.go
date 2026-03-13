package controller

import (
	"net/http"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/internal/middlewares"
	"github.com/BangNopall/paskihub-be/pkg/helpers/http/response"
	"github.com/BangNopall/paskihub-be/pkg/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type eventController struct {
	eventSvc contracts.EventService
	redis    redis.RedisInterface
}

func InitEventController(
	eventSvc contracts.EventService,
	router fiber.Router,
	middleware *middlewares.Middleware,
	redis redis.RedisInterface,
) {
	eventController := &eventController{
		eventSvc: eventSvc,
		redis:    redis,
	}

	eventRouter := router.Group("/api/v1/events")

	eventRouter.Post("/create", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, eventController.CreateEvent)
	eventRouter.Get("/user/:userId", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, eventController.ShowUserEvent)
	eventRouter.Post("/upload/:id/logo", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, eventController.UploadLogo)
	eventRouter.Post("/upload/:id/poster", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, eventController.UploadPoster)
	eventRouter.Put("/update/:id", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, eventController.UpdateEvent)
	eventRouter.Delete("/delete/:id", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, eventController.DeleteEvent)

	// admin
	eventRouter.Get("/admin/show/:id", middleware.Authentication, middleware.RateLimiter(), middleware.AuthAdmin, eventController.ShowEventData)

	// Event Levels
	eventRouter.Post("/:id/levels", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, eventController.CreateEventLevel)
	eventRouter.Put("/:id/levels/:levelId", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, eventController.UpdateEventLevel)
	eventRouter.Delete("/:id/levels/:levelId", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, eventController.DeleteEventLevel)
}

// CreateEvent godoc
// @Summary Create a new event
// @Description Create an event for the organizer
// @Tags Events
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param event body dto.EventCreate true "Event details"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/events/create [post]
func (c *eventController) CreateEvent(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to create event"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	var event dto.EventCreate
	if err = ctx.BodyParser(&event); err != nil {
		return nil
	}

	userIdStr := ctx.Locals("id").(string)
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		code = http.StatusUnauthorized
		return nil
	}

	err = c.eventSvc.CreateEvent(ctx.Context(), userId, event)
	code = domain.GetCode(err)
	if err != nil {
		return nil
	}

	message = "success to create event"
	return nil
}

// ShowEventData godoc
// @Summary Show event data (Admin)
// @Description Get specific event details by ID
// @Tags Events
// @Security BearerAuth
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/events/admin/show/{id} [get]
func (c *eventController) ShowEventData(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to show event data"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	idStr := ctx.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil
	}

	res, err = c.eventSvc.ShowEventData(ctx.Context(), id)
	code = domain.GetCode(err)
	if err != nil {
		return nil
	}

	message = "success to show event data"
	return nil
}

// ShowUserEvent godoc
// @Summary Show auth user's events
// @Description Get all events created by specific user/organizer
// @Tags Events
// @Security BearerAuth
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/events/user/{userId} [get]
func (c *eventController) ShowUserEvent(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to show user events"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	userIdStr := ctx.Params("userId")
	authUserIdStr := ctx.Locals("id").(string)
	if authUserIdStr != userIdStr {
		code = http.StatusForbidden
		message = "forbidden access"
		err = domain.ErrForbidden
		return nil
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return nil
	}

	res, err = c.eventSvc.ShowUserEvent(ctx.Context(), userId)
	code = domain.GetCode(err)
	if err != nil {
		return nil
	}

	message = "success to show user events"
	return nil
}

// UploadLogo godoc
// @Summary Upload event logo
// @Description Upload logo image for an event
// @Tags Events
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Event ID"
// @Param logo formData file true "Logo Image"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/events/upload/{id}/logo [post]
func (c *eventController) UploadLogo(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to upload logo"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	id := ctx.Params("id")
	userId := ctx.Locals("id").(string)

	file, err := ctx.FormFile("logo")
	if err != nil {
		return nil
	}

	err = c.eventSvc.UploadLogo(ctx.Context(), id, userId, file)
	code = domain.GetCode(err)
	if err != nil {
		return nil
	}

	message = "success to upload logo"
	return nil
}

// UploadPoster godoc
// @Summary Upload event poster
// @Description Upload poster image for an event
// @Tags Events
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Event ID"
// @Param poster formData file true "Poster Image"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/events/upload/{id}/poster [post]
func (c *eventController) UploadPoster(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to upload poster"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	id := ctx.Params("id")
	userId := ctx.Locals("id").(string)

	file, err := ctx.FormFile("poster")
	if err != nil {
		return nil
	}

	err = c.eventSvc.UploadPoster(ctx.Context(), id, userId, file)
	code = domain.GetCode(err)
	if err != nil {
		return nil
	}

	message = "success to upload poster"
	return nil
}

// UpdateEvent godoc
// @Summary Update an event
// @Description Update event details
// @Tags Events
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param event body dto.EventUpdate true "Updated details"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/events/update/{id} [put]
func (c *eventController) UpdateEvent(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to update event"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	id := ctx.Params("id")
	userId := ctx.Locals("id").(string)

	var event dto.EventUpdate
	if err = ctx.BodyParser(&event); err != nil {
		return nil
	}

	err = c.eventSvc.UpdateEvent(ctx.Context(), id, userId, &event)
	code = domain.GetCode(err)
	if err != nil {
		return nil
	}

	message = "success to update event"
	return nil
}

// DeleteEvent godoc
// @Summary Delete an event
// @Description Remove an event by ID
// @Tags Events
// @Security BearerAuth
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/events/delete/{id} [delete]
func (c *eventController) DeleteEvent(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to delete event"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	id := ctx.Params("id")
	userId := ctx.Locals("id").(string)

	err = c.eventSvc.DeleteEvent(ctx.Context(), id, userId)
	code = domain.GetCode(err)
	if err != nil {
		return nil
	}

	message = "success to delete event"
	return nil
}

// CreateEventLevel godoc
// @Summary Create event level
// @Description Create a new level category under an event
// @Tags Events
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param level body dto.EventLevelCreate true "Level details"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/events/{id}/levels [post]
func (c *eventController) CreateEventLevel(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to create event level"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	id := ctx.Params("id")
	userId := ctx.Locals("id").(string)

	var level dto.EventLevelCreate
	if err = ctx.BodyParser(&level); err != nil {
		return nil
	}

	err = c.eventSvc.CreateEventLevel(ctx.Context(), id, userId, &level)
	code = domain.GetCode(err)
	if err != nil {
		return nil
	}

	message = "success to create event level"
	return nil
}

// UpdateEventLevel godoc
// @Summary Update event level
// @Description Update a level's detail under an event
// @Tags Events
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param levelId path string true "Level ID"
// @Param level body dto.EventLevelUpdate true "Updated level details"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/events/{id}/levels/{levelId} [put]
func (c *eventController) UpdateEventLevel(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to update event level"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	id := ctx.Params("id")
	levelId := ctx.Params("levelId")
	userId := ctx.Locals("id").(string)

	levelUUID, err := uuid.Parse(levelId)
	if err != nil {
		return nil
	}

	eventId, err := uuid.Parse(id)
	if err != nil {
		return nil
	}

	var level dto.EventLevelUpdate
	if err = ctx.BodyParser(&level); err != nil {
		return nil
	}

	// Ensure IDs match params
	level.Id = levelUUID
	level.EventId = eventId

	err = c.eventSvc.UpdateEventLevel(ctx.Context(), id, userId, &level)
	code = domain.GetCode(err)
	if err != nil {
		return nil
	}

	message = "success to update event level"
	return nil
}

// DeleteEventLevel godoc
// @Summary Delete event level
// @Description Remove a level from an event
// @Tags Events
// @Security BearerAuth
// @Produce json
// @Param id path string true "Event ID"
// @Param levelId path string true "Level ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/events/{id}/levels/{levelId} [delete]
func (c *eventController) DeleteEventLevel(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to delete event level"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	id := ctx.Params("id")
	levelId := ctx.Params("levelId")
	userId := ctx.Locals("id").(string)

	err = c.eventSvc.DeleteEventLevel(ctx.Context(), id, levelId, userId)
	code = domain.GetCode(err)
	if err != nil {
		return nil
	}

	message = "success to delete event level"
	return nil
}
