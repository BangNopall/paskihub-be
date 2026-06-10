package controller

import (
	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/internal/middlewares"
	"github.com/BangNopall/paskihub-be/pkg/helpers/http/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type privateFileController struct {
	service contracts.PrivateFileService
}

func InitPrivateFileController(
	service contracts.PrivateFileService,
	router fiber.Router,
	middleware *middlewares.Middleware,
) {
	controller := &privateFileController{service: service}
	router.Get(
		"/api/v1/files/:resourceType/:resourceId",
		middleware.Authentication,
		middleware.RateLimiter(),
		controller.Download,
	)
}

// Download godoc
// @Summary Download an authorized private file
// @Description Streams a sensitive file after resource ownership validation
// @Tags Private Files
// @Security ApiKeyAuth && BearerAuth
// @Produce application/octet-stream
// @Param resourceType path string true "Private resource type"
// @Param resourceId path string true "Resource ID"
// @Success 200 {file} binary
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/files/{resourceType}/{resourceId} [get]
func (c *privateFileController) Download(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Locals("id").(string))
	if err != nil {
		response.Send(ctx, fiber.StatusUnauthorized, "failed to authenticate user", nil, err)
		return nil
	}

	if parentID, ok := ctx.Locals("parent_id").(string); ok && parentID != "" {
		if parsedParentID, parseErr := uuid.Parse(parentID); parseErr == nil {
			userID = parsedParentID
		}
	}

	resourceID, err := uuid.Parse(ctx.Params("resourceId"))
	if err != nil {
		response.Send(ctx, fiber.StatusBadRequest, "invalid resource id", nil, err)
		return nil
	}

	path, err := c.service.ResolvePrivateFile(
		ctx.Context(),
		userID,
		ctx.Locals("role").(string),
		ctx.Params("resourceType"),
		resourceID,
	)
	if err != nil {
		response.Send(ctx, domain.GetCode(err), "failed to download private file", nil, err)
		return nil
	}

	return ctx.SendFile(path)
}
