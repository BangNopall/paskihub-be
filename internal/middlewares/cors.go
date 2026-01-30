package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:"http://localhost:3000",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Content-Type, X-XSRF-TOKEN, Accept, Origin, X-Requested-With, Authorization, X-API-Key, X-Cursor, Token-Type",
		ExposeHeaders:"Content-Length",
		AllowCredentials: true,
	})
}
