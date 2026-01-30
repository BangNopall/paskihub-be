package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/BangNopall/paskihub-be/pkg/helpers/http/response"
	"github.com/BangNopall/paskihub-be/pkg/log"
	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) Authentication(ctx *fiber.Ctx) error {
	bearer := ctx.Get("Authorization")
	if bearer == "" {
		log.Warn(log.LogInfo{
			"error": errors.New("failed to get bearer token"),
		}, "[MIDDLEWARE][Authentication] failed to get bearer token")

		response.SendErrResp(
			ctx,
			http.StatusUnauthorized,
			response.Error,
			"failed to authenticate user",
			errors.New("failed to get bearer token"),
		)
		return nil
	}

	splitted := strings.Split(bearer, " ")

	if len(splitted) < 2 {
		response.SendErrResp(
			ctx,
			400,
			response.Fail,
			"failed to authenticate user",
			fmt.Errorf("invalid token"),
		)
		return nil
	}

	tokenString := splitted[1]
	id, email, role, err := m.jwt.ValidateToken(tokenString)
	if err != nil {
		log.Warn(log.LogInfo{
			"error": err,
		}, "[MIDDLEWARE][Authentication] failed to validate token")

		response.SendErrResp(
			ctx,
			http.StatusUnauthorized,
			response.Error,
			"failed to authenticate user",
			err,
		)
		return nil
	}

	val, err := m.redis.Get(ctx.Context(), tokenString)

	if err != nil {
		response.SendErrResp(
			ctx,
			http.StatusInternalServerError,
			response.Error,
			"failed to authenticate user",
			err,
		)
		return nil
	}

	if val != "" {
		response.SendErrResp(
			ctx,
			http.StatusUnauthorized,
			response.Fail,
			"failed to authenticate user",
			nil,
		)
		return nil
	}

	ctx.Locals("id", id.String())
	ctx.Locals("email", email)
	ctx.Locals("role", role)
	ctx.Next()
	return nil
}