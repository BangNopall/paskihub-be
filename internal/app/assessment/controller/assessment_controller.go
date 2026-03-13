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

// CreateJudge godoc
// @Summary Create judge
// @Description Create a judge for an event
// @Tags Judges
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Param req body dto.CreateJudgeReq true "Judge Details"
// @Success 201 {object} map[string]interface{}
// @Router /api/v1/eo/events/{eventId}/assessment/judges [post]
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

// GetJudges godoc
// @Summary Get judges
// @Description Get all judges for an event
// @Tags Judges
// @Security BearerAuth
// @Produce json
// @Param eventId path string true "Event ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/eo/events/{eventId}/assessment/judges [get]
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

// UpdateJudge godoc
// @Summary Update judge
// @Description Update judge details
// @Tags Judges
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Param id path string true "Judge ID"
// @Param req body dto.UpdateJudgeReq true "Updated Judge Details"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/eo/events/{eventId}/assessment/judges/{id} [put]
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

// DeleteJudge godoc
// @Summary Delete judge
// @Description Remove judge from an event
// @Tags Judges
// @Security BearerAuth
// @Produce json
// @Param eventId path string true "Event ID"
// @Param id path string true "Judge ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/eo/events/{eventId}/assessment/judges/{id} [delete]
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

// CreateViolationType godoc
// @Summary Create violation type
// @Description Create a new violation type
// @Tags Violation Types
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Param req body dto.CreateViolationTypeReq true "Violation Type Details"
// @Success 201 {object} map[string]interface{}
// @Router /api/v1/eo/events/{eventId}/assessment/violation-types [post]
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

// GetViolationTypes godoc
// @Summary Get violation types
// @Description Get all violation types for an event
// @Tags Violation Types
// @Security BearerAuth
// @Produce json
// @Param eventId path string true "Event ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/eo/events/{eventId}/assessment/violation-types [get]
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

// UpdateViolationType godoc
// @Summary Update violation type
// @Description Update violation type details
// @Tags Violation Types
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Param id path string true "Violation Type ID"
// @Param req body dto.UpdateViolationTypeReq true "Updated Violation Type Details"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/eo/events/{eventId}/assessment/violation-types/{id} [put]
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

// DeleteViolationType godoc
// @Summary Delete violation type
// @Description Remove violation type from an event
// @Tags Violation Types
// @Security BearerAuth
// @Produce json
// @Param eventId path string true "Event ID"
// @Param id path string true "Violation Type ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/eo/events/{eventId}/assessment/violation-types/{id} [delete]
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

// CreateScoreCategory godoc
// @Summary Create score category
// @Description Create a new score category
// @Tags Score Categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Param req body dto.CreateScoreCategoryReq true "Score Category Details"
// @Success 201 {object} map[string]interface{}
// @Router /api/v1/eo/events/{eventId}/assessment/score-categories [post]
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

// GetScoreCategories godoc
// @Summary Get score categories
// @Description Get all score categories for an event
// @Tags Score Categories
// @Security BearerAuth
// @Produce json
// @Param eventId path string true "Event ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/eo/events/{eventId}/assessment/score-categories [get]
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

// UpdateScoreCategory godoc
// @Summary Update score category
// @Description Update score category details
// @Tags Score Categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Param id path string true "Score Category ID"
// @Param req body dto.UpdateScoreCategoryReq true "Updated Score Category Details"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/eo/events/{eventId}/assessment/score-categories/{id} [put]
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

// DeleteScoreCategory godoc
// @Summary Delete score category
// @Description Remove score category from an event
// @Tags Score Categories
// @Security BearerAuth
// @Produce json
// @Param eventId path string true "Event ID"
// @Param id path string true "Score Category ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/eo/events/{eventId}/assessment/score-categories/{id} [delete]
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

// CreateScoreSubCategory godoc
// @Summary Create score sub category
// @Description Create a new score sub category
// @Tags Score Sub Categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Param req body dto.CreateScoreSubCategoryReq true "Score Sub Category Details"
// @Success 201 {object} map[string]interface{}
// @Router /api/v1/eo/events/{eventId}/assessment/score-sub-categories [post]
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

// UpdateScoreSubCategory godoc
// @Summary Update score sub category
// @Description Update score sub category details
// @Tags Score Sub Categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Param id path string true "Score Sub Category ID"
// @Param req body dto.UpdateScoreSubCategoryReq true "Updated Score Sub Category Details"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/eo/events/{eventId}/assessment/score-sub-categories/{id} [put]
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

// DeleteScoreSubCategory godoc
// @Summary Delete score sub category
// @Description Remove score sub category from an event
// @Tags Score Sub Categories
// @Security BearerAuth
// @Produce json
// @Param eventId path string true "Event ID"
// @Param id path string true "Score Sub Category ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/eo/events/{eventId}/assessment/score-sub-categories/{id} [delete]
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

// SetupGradeRules godoc
// @Summary Setup grade rules
// @Description Configure grade rule mapping for an event
// @Tags Assessment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Param req body dto.SetupGradeRulesReq true "Grade Rules Details"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/eo/events/{eventId}/assessment/grade-rules [post]
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

// GetGradeRules godoc
// @Summary Get grade rules
// @Description Get existing grade rules for an event
// @Tags Assessment
// @Security BearerAuth
// @Produce json
// @Param eventId path string true "Event ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/eo/events/{eventId}/assessment/grade-rules [get]
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

// InputScore godoc
// @Summary Input single score
// @Description Input an individual score
// @Tags Assessment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Param req body dto.InputScoreReq true "Score Input Details"
// @Success 201 {object} map[string]interface{}
// @Router /api/v1/eo/events/{eventId}/assessment/scores [post]
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
