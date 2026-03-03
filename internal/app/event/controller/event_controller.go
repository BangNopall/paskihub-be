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
