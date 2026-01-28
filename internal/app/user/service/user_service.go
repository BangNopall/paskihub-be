package service

import (
	"context"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/contracts"
	"github.com/hology8/hology-be/domain/dto"
	"github.com/hology8/hology-be/domain/entity"
	"github.com/hology8/hology-be/internal/infra/env"
	"github.com/hology8/hology-be/pkg/aws"
	"github.com/hology8/hology-be/pkg/bcrypt"
	"github.com/hology8/hology-be/pkg/gomail"
	"github.com/hology8/hology-be/pkg/helpers"
	html_content "github.com/hology8/hology-be/pkg/html"
	"github.com/hology8/hology-be/pkg/jwt"
	"github.com/hology8/hology-be/pkg/log"
	"github.com/hology8/hology-be/pkg/redis"
	timePkg "github.com/hology8/hology-be/pkg/time"
	uuidPkg "github.com/hology8/hology-be/pkg/uuid"
)

type userService struct {
	userRepo contracts.UserRepository
	uuid     uuidPkg.UUIDInterface
	bcrypt   bcrypt.BcryptInterface
	time     timePkg.TimeInterface
	goMail   gomail.GoMailInterface
	jwt      jwt.JwtInterface
	redis    redis.RedisInterface
	timeout  time.Duration
}

func NewUserService(
	userRepo contracts.UserRepository,
	uuid uuidPkg.UUIDInterface,
	bcrypt bcrypt.BcryptInterface,
	time timePkg.TimeInterface,
	goMail gomail.GoMailInterface,
	jwt jwt.JwtInterface,
	redis redis.RedisInterface,
	timeout time.Duration,
) contracts.UserService {
	return &userService{
		userRepo,
		uuid,
		bcrypt,
		time,
		goMail,
		jwt,
		redis,
		timeout,
	}
}