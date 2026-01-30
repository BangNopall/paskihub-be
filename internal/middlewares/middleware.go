package middlewares

import (
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/pkg/jwt"
	"github.com/BangNopall/paskihub-be/pkg/redis"
)

type Middleware struct {
	jwt       jwt.JwtInterface
	userRepo  contracts.UserRepository
	redis     redis.RedisInterface
}

func NewMiddleware(
	jwt jwt.JwtInterface,
	userRepo contracts.UserRepository,
	redis redis.RedisInterface,
) *Middleware {
	return &Middleware{
		jwt:       jwt,
		userRepo:  userRepo,
		redis:     redis,
	}
}