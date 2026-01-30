package middlewares

import (
	"errors"
	"time"

	"github.com/BangNopall/paskihub-be/pkg/helpers/http/response"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func (m *Middleware) RateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        60,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			response.SendErrResp(
				c,
				fiber.StatusTooManyRequests,
				response.Fail,
				"too many requests",
				errors.New("too many requests, please try again later"),
			)
			return nil
		},
	})
}
