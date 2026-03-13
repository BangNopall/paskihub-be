package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/BangNopall/paskihub-be/internal/infra/env"
	"github.com/BangNopall/paskihub-be/pkg/helpers/http/response"
)

func ApiKey() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Bypass API key check for Swagger UI paths
		if strings.HasPrefix(ctx.Path(), "/swagger") {
			return ctx.Next()
		}

		headerReq := ctx.Get("x-api-key")

		if headerReq == "" {
			response.SendErrResp(
				ctx,
				http.StatusBadRequest,
				response.Fail,
				"failed to authenticate request",
				fmt.Errorf("invalid api key"),
			)
			return nil
		}

		splitted := strings.Split(headerReq, " ")
		if len(splitted) != 2 || splitted[0] != "Key" {
			response.SendErrResp(
				ctx,
				http.StatusBadRequest,
				response.Fail,
				"failed to authenticate request",
				fmt.Errorf("invalid api key format, use: Key <api_key>"),
			)
			return nil
		}
		headerKey := splitted[1]

		if headerKey != env.AppEnv.ApiKey {
			response.SendErrResp(
				ctx,
				http.StatusBadRequest,
				response.Fail,
				"failed to authenticate request",
				fmt.Errorf("invalid api key"),
			)
			return nil
		}

		ctx.Next()
		return nil
	}
}
