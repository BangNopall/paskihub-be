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
		headerReq := ctx.Get("x-api-key")

		splitted := strings.Split(headerReq, " ")

		if len(splitted) < 2 {
			response.SendErrResp(
				ctx,
				http.StatusBadRequest,
				response.Fail,
				"failed to authenticate request",
				fmt.Errorf("invalid api key"),
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