package controller

import (
	"net/http"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/internal/middlewares"
	"github.com/BangNopall/paskihub-be/pkg/helpers/http/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type assessmentController struct {
	svc contracts.IAssessmentService
}

func InitAssessmentController(
	svc contracts.IAssessmentService,
	router fiber.Router,
	middleware *middlewares.Middleware,
) {
	c := &assessmentController{svc: svc}
	
	group := router.Group("/api/v1/eo/events/:eventId/assessment")

	// Judges
	group.Post("/judges", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.CreateJudge)
	group.Get("/judges", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.GetJudges)
	group.Put("/judges/:id", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.UpdateJudge)
	group.Delete("/judges/:id", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.DeleteJudge)

	// Violation Types
	group.Post("/violation-types", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.CreateViolationType)
	group.Get("/violation-types", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.GetViolationTypes)
	group.Put("/violation-types/:id", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.UpdateViolationType)
	group.Delete("/violation-types/:id", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.DeleteViolationType)

	// Score Categories
	group.Post("/score-categories", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.CreateScoreCategory)
	group.Get("/score-categories", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.GetScoreCategories)
	group.Put("/score-categories/:id", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.UpdateScoreCategory)
	group.Delete("/score-categories/:id", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.DeleteScoreCategory)

	// Score Sub Categories
	group.Post("/score-sub-categories", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.CreateScoreSubCategory)
	group.Put("/score-sub-categories/:id", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.UpdateScoreSubCategory)
	group.Delete("/score-sub-categories/:id", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.DeleteScoreSubCategory)

	// Grade Rules
	group.Post("/grade-rules", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.SetupGradeRules)
	group.Get("/grade-rules", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.GetGradeRules)

	// Score Input
	group.Post("/scores", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.InputScore)
}

func getUUIDParam(ctx *fiber.Ctx, key string) (uuid.UUID, error) {
	val := ctx.Params(key)
	return uuid.Parse(val)
}

func getUserId(ctx *fiber.Ctx) (uuid.UUID, error) {
	userIdStr := ctx.Locals("id").(string)
	return uuid.Parse(userIdStr)
}

// Judge Handlers
func (c *assessmentController) CreateJudge(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusCreated
		res     interface{}
		message string = "failed to create judge"
	)
	sendResp := func() { response.Send(ctx, code, message, res, err) }
	defer sendResp()

	userId, err := getUserId(ctx)
	if err != nil { code = http.StatusUnauthorized; return nil }
	eventId, err := getUUIDParam(ctx, "eventId")
	if err != nil { code = http.StatusBadRequest; return nil }

	var req dto.CreateJudgeReq
	if errParse := ctx.BodyParser(&req); errParse != nil {
		err = errParse
		code = http.StatusBadRequest
		return nil
	}

	res, err = c.svc.CreateJudge(ctx.Context(), eventId, userId, req)
	code = domain.GetCode(err)
	if err == nil {
		message = "success to create judge"
		code = http.StatusCreated
	}
	return nil
}

func (c *assessmentController) GetJudges(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusOK
		res     interface{}
		message string = "failed to get judges"
	)
	sendResp := func() { response.Send(ctx, code, message, res, err) }
	defer sendResp()

	userId, err := getUserId(ctx)
	if err != nil { code = http.StatusUnauthorized; return nil }
	eventId, err := getUUIDParam(ctx, "eventId")
	if err != nil { code = http.StatusBadRequest; return nil }

	res, err = c.svc.GetJudges(ctx.Context(), eventId, userId)
	code = domain.GetCode(err)
	if err == nil { message = "success to get judges" }
	return nil
}

func (c *assessmentController) UpdateJudge(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusOK
		res     interface{}
		message string = "failed to update judge"
	)
	sendResp := func() { response.Send(ctx, code, message, res, err) }
	defer sendResp()

	userId, err := getUserId(ctx)
	if err != nil { code = http.StatusUnauthorized; return nil }
	eventId, err := getUUIDParam(ctx, "eventId")
	if err != nil { code = http.StatusBadRequest; return nil }
	id, err := getUUIDParam(ctx, "id")
	if err != nil { code = http.StatusBadRequest; return nil }

	var req dto.UpdateJudgeReq
	if errParse := ctx.BodyParser(&req); errParse != nil {
		err = errParse; code = http.StatusBadRequest; return nil
	}

	res, err = c.svc.UpdateJudge(ctx.Context(), eventId, userId, id, req)
	code = domain.GetCode(err)
	if err == nil { message = "success to update judge" }
	return nil
}

func (c *assessmentController) DeleteJudge(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusOK
		message string = "failed to delete judge"
	)
	sendResp := func() { response.Send(ctx, code, message, nil, err) }
	defer sendResp()

	userId, err := getUserId(ctx)
	if err != nil { code = http.StatusUnauthorized; return nil }
	eventId, err := getUUIDParam(ctx, "eventId")
	if err != nil { code = http.StatusBadRequest; return nil }
	id, err := getUUIDParam(ctx, "id")
	if err != nil { code = http.StatusBadRequest; return nil }

	err = c.svc.DeleteJudge(ctx.Context(), eventId, userId, id)
	code = domain.GetCode(err)
	if err == nil { message = "success to delete judge" }
	return nil
}

// ViolationType Handlers
func (c *assessmentController) CreateViolationType(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusCreated
		res     interface{}
		message string = "failed to create violation type"
	)
	sendResp := func() { response.Send(ctx, code, message, res, err) }
	defer sendResp()

	userId, err := getUserId(ctx)
	if err != nil { code = http.StatusUnauthorized; return nil }
	eventId, err := getUUIDParam(ctx, "eventId")
	if err != nil { code = http.StatusBadRequest; return nil }

	var req dto.CreateViolationTypeReq
	if errParse := ctx.BodyParser(&req); errParse != nil {
		err = errParse; code = http.StatusBadRequest; return nil
	}

	res, err = c.svc.CreateViolationType(ctx.Context(), eventId, userId, req)
	code = domain.GetCode(err)
	if err == nil {
		message = "success to create violation type"
		code = http.StatusCreated
	}
	return nil
}

func (c *assessmentController) GetViolationTypes(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusOK
		res     interface{}
		message string = "failed to get violation types"
	)
	sendResp := func() { response.Send(ctx, code, message, res, err) }
	defer sendResp()

	userId, err := getUserId(ctx)
	if err != nil { code = http.StatusUnauthorized; return nil }
	eventId, err := getUUIDParam(ctx, "eventId")
	if err != nil { code = http.StatusBadRequest; return nil }

	res, err = c.svc.GetViolationTypes(ctx.Context(), eventId, userId)
	code = domain.GetCode(err)
	if err == nil { message = "success to get violation types" }
	return nil
}

func (c *assessmentController) UpdateViolationType(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusOK
		res     interface{}
		message string = "failed to update violation type"
	)
	sendResp := func() { response.Send(ctx, code, message, res, err) }
	defer sendResp()

	userId, err := getUserId(ctx)
	if err != nil { code = http.StatusUnauthorized; return nil }
	eventId, err := getUUIDParam(ctx, "eventId")
	if err != nil { code = http.StatusBadRequest; return nil }
	id, err := getUUIDParam(ctx, "id")
	if err != nil { code = http.StatusBadRequest; return nil }

	var req dto.UpdateViolationTypeReq
	if errParse := ctx.BodyParser(&req); errParse != nil {
		err = errParse; code = http.StatusBadRequest; return nil
	}

	res, err = c.svc.UpdateViolationType(ctx.Context(), eventId, userId, id, req)
	code = domain.GetCode(err)
	if err == nil { message = "success to update violation type" }
	return nil
}

func (c *assessmentController) DeleteViolationType(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusOK
		message string = "failed to delete violation type"
	)
	sendResp := func() { response.Send(ctx, code, message, nil, err) }
	defer sendResp()

	userId, err := getUserId(ctx)
	if err != nil { code = http.StatusUnauthorized; return nil }
	eventId, err := getUUIDParam(ctx, "eventId")
	if err != nil { code = http.StatusBadRequest; return nil }
	id, err := getUUIDParam(ctx, "id")
	if err != nil { code = http.StatusBadRequest; return nil }

	err = c.svc.DeleteViolationType(ctx.Context(), eventId, userId, id)
	code = domain.GetCode(err)
	if err == nil { message = "success to delete violation type" }
	return nil
}

// ScoreCategory Handlers
func (c *assessmentController) CreateScoreCategory(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusCreated
		res     interface{}
		message string = "failed to create score category"
	)
	sendResp := func() { response.Send(ctx, code, message, res, err) }
	defer sendResp()

	userId, err := getUserId(ctx)
	if err != nil { code = http.StatusUnauthorized; return nil }
	eventId, err := getUUIDParam(ctx, "eventId")
	if err != nil { code = http.StatusBadRequest; return nil }

	var req dto.CreateScoreCategoryReq
	if errParse := ctx.BodyParser(&req); errParse != nil {
		err = errParse; code = http.StatusBadRequest; return nil
	}

	res, err = c.svc.CreateScoreCategory(ctx.Context(), eventId, userId, req)
	code = domain.GetCode(err)
	if err == nil {
		message = "success to create score category"
		code = http.StatusCreated
	}
	return nil
}

func (c *assessmentController) GetScoreCategories(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusOK
		res     interface{}
		message string = "failed to get score categories"
	)
	sendResp := func() { response.Send(ctx, code, message, res, err) }
	defer sendResp()

	userId, err := getUserId(ctx)
	if err != nil { code = http.StatusUnauthorized; return nil }
	eventId, err := getUUIDParam(ctx, "eventId")
	if err != nil { code = http.StatusBadRequest; return nil }

	res, err = c.svc.GetScoreCategories(ctx.Context(), eventId, userId)
	code = domain.GetCode(err)
	if err == nil { message = "success to get score categories" }
	return nil
}

func (c *assessmentController) UpdateScoreCategory(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusOK
		res     interface{}
		message string = "failed to update score category"
	)
	sendResp := func() { response.Send(ctx, code, message, res, err) }
	defer sendResp()

	userId, err := getUserId(ctx)
	if err != nil { code = http.StatusUnauthorized; return nil }
	eventId, err := getUUIDParam(ctx, "eventId")
	if err != nil { code = http.StatusBadRequest; return nil }
	id, err := getUUIDParam(ctx, "id")
	if err != nil { code = http.StatusBadRequest; return nil }

	var req dto.UpdateScoreCategoryReq
	if errParse := ctx.BodyParser(&req); errParse != nil {
		err = errParse; code = http.StatusBadRequest; return nil
	}

	res, err = c.svc.UpdateScoreCategory(ctx.Context(), eventId, userId, id, req)
	code = domain.GetCode(err)
	if err == nil { message = "success to update score category" }
	return nil
}

func (c *assessmentController) DeleteScoreCategory(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusOK
		message string = "failed to delete score category"
	)
	sendResp := func() { response.Send(ctx, code, message, nil, err) }
	defer sendResp()

	userId, err := getUserId(ctx)
	if err != nil { code = http.StatusUnauthorized; return nil }
	eventId, err := getUUIDParam(ctx, "eventId")
	if err != nil { code = http.StatusBadRequest; return nil }
	id, err := getUUIDParam(ctx, "id")
	if err != nil { code = http.StatusBadRequest; return nil }

	err = c.svc.DeleteScoreCategory(ctx.Context(), eventId, userId, id)
	code = domain.GetCode(err)
	if err == nil { message = "success to delete score category" }
	return nil
}

// ScoreSubCategory Handlers
func (c *assessmentController) CreateScoreSubCategory(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusCreated
		res     interface{}
		message string = "failed to create score sub category"
	)
	sendResp := func() { response.Send(ctx, code, message, res, err) }
	defer sendResp()

	userId, err := getUserId(ctx)
	if err != nil { code = http.StatusUnauthorized; return nil }
	eventId, err := getUUIDParam(ctx, "eventId")
	if err != nil { code = http.StatusBadRequest; return nil }

	var req dto.CreateScoreSubCategoryReq
	if errParse := ctx.BodyParser(&req); errParse != nil {
		err = errParse; code = http.StatusBadRequest; return nil
	}

	res, err = c.svc.CreateScoreSubCategory(ctx.Context(), eventId, userId, req)
	code = domain.GetCode(err)
	if err == nil {
		message = "success to create score sub category"
		code = http.StatusCreated
	}
	return nil
}

func (c *assessmentController) UpdateScoreSubCategory(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusOK
		res     interface{}
		message string = "failed to update score sub category"
	)
	sendResp := func() { response.Send(ctx, code, message, res, err) }
	defer sendResp()

	userId, err := getUserId(ctx)
	if err != nil { code = http.StatusUnauthorized; return nil }
	eventId, err := getUUIDParam(ctx, "eventId")
	if err != nil { code = http.StatusBadRequest; return nil }
	id, err := getUUIDParam(ctx, "id")
	if err != nil { code = http.StatusBadRequest; return nil }

	var req dto.UpdateScoreSubCategoryReq
	if errParse := ctx.BodyParser(&req); errParse != nil {
		err = errParse; code = http.StatusBadRequest; return nil
	}

	res, err = c.svc.UpdateScoreSubCategory(ctx.Context(), eventId, userId, id, req)
	code = domain.GetCode(err)
	if err == nil { message = "success to update score sub category" }
	return nil
}

func (c *assessmentController) DeleteScoreSubCategory(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusOK
		message string = "failed to delete score sub category"
	)
	sendResp := func() { response.Send(ctx, code, message, nil, err) }
	defer sendResp()

	userId, err := getUserId(ctx)
	if err != nil { code = http.StatusUnauthorized; return nil }
	eventId, err := getUUIDParam(ctx, "eventId")
	if err != nil { code = http.StatusBadRequest; return nil }
	id, err := getUUIDParam(ctx, "id")
	if err != nil { code = http.StatusBadRequest; return nil }

	err = c.svc.DeleteScoreSubCategory(ctx.Context(), eventId, userId, id)
	code = domain.GetCode(err)
	if err == nil { message = "success to delete score sub category" }
	return nil
}

// GradeRules Handlers
func (c *assessmentController) SetupGradeRules(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusOK
		res     interface{}
		message string = "failed to setup grade rules"
	)
	sendResp := func() { response.Send(ctx, code, message, res, err) }
	defer sendResp()

	userId, err := getUserId(ctx)
	if err != nil { code = http.StatusUnauthorized; return nil }
	eventId, err := getUUIDParam(ctx, "eventId")
	if err != nil { code = http.StatusBadRequest; return nil }

	var req dto.SetupGradeRulesReq
	if errParse := ctx.BodyParser(&req); errParse != nil {
		err = errParse; code = http.StatusBadRequest; return nil
	}

	res, err = c.svc.SetupGradeRules(ctx.Context(), eventId, userId, req)
	code = domain.GetCode(err)
	if err == nil { message = "success to setup grade rules" }
	return nil
}

func (c *assessmentController) GetGradeRules(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusOK
		res     interface{}
		message string = "failed to get grade rules"
	)
	sendResp := func() { response.Send(ctx, code, message, res, err) }
	defer sendResp()

	userId, err := getUserId(ctx)
	if err != nil { code = http.StatusUnauthorized; return nil }
	eventId, err := getUUIDParam(ctx, "eventId")
	if err != nil { code = http.StatusBadRequest; return nil }

	res, err = c.svc.GetGradeRules(ctx.Context(), eventId, userId)
	code = domain.GetCode(err)
	if err == nil { message = "success to get grade rules" }
	return nil
}

// Score Handlers
func (c *assessmentController) InputScore(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusCreated
		res     interface{}
		message string = "failed to input score"
	)
	sendResp := func() { response.Send(ctx, code, message, res, err) }
	defer sendResp()

	userId, err := getUserId(ctx)
	if err != nil { code = http.StatusUnauthorized; return nil }
	eventId, err := getUUIDParam(ctx, "eventId")
	if err != nil { code = http.StatusBadRequest; return nil }

	var req dto.InputScoreReq
	if errParse := ctx.BodyParser(&req); errParse != nil {
		err = errParse; code = http.StatusBadRequest; return nil
	}

	res, err = c.svc.InputScore(ctx.Context(), eventId, userId, req)
	code = domain.GetCode(err)
	if err == nil {
		message = "success to input score"
		code = http.StatusCreated
	}
	return nil
}
