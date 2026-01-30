package middlewares

import (
	"net/http"

	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/pkg/helpers/http/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (m *Middleware) Verified(ctx *fiber.Ctx) error {
	id := ctx.Locals("id")
	if id == nil {
		response.SendErrResp(
			ctx,
			http.StatusUnauthorized,
			response.Error,
			"unauthorized",
			nil,
		)
		return nil
	}

	userIdStr, ok := id.(string)
	if !ok {
		response.SendErrResp(
			ctx,
			http.StatusInternalServerError,
			response.Error,
			"invalid user id type",
			nil,
		)
		return nil
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		response.SendErrResp(
			ctx,
			http.StatusUnauthorized,
			response.Error,
			"invalid user id",
			err,
		)
		return nil
	}

	var user entity.User
	err = m.userRepo.FindUser(&user, &dto.UserParam{ID: userId})
	if err != nil {
		response.SendErrResp(
			ctx,
			http.StatusInternalServerError,
			response.Error,
			"user not found",
			err,
		)
		return nil
	}

	if !user.EmailIsVerified {
		response.SendErrResp(
			ctx,
			http.StatusForbidden,
			response.Fail,
			"email not verified",
			nil,
		)
		return nil
	}

	return ctx.Next()
}
