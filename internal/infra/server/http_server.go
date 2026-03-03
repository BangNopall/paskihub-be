package server

import (
	"time"

	"github.com/BangNopall/paskihub-be/internal/middlewares"
	"github.com/BangNopall/paskihub-be/pkg/bcrypt"
	validators "github.com/BangNopall/paskihub-be/pkg/fiber"
	"github.com/BangNopall/paskihub-be/pkg/gomail"
	"github.com/BangNopall/paskihub-be/pkg/jwt"
	"github.com/BangNopall/paskihub-be/pkg/log"
	"github.com/BangNopall/paskihub-be/pkg/redis"
	timePkg "github.com/BangNopall/paskihub-be/pkg/time"
	"github.com/BangNopall/paskihub-be/pkg/uuid"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"

	userCtr "github.com/BangNopall/paskihub-be/internal/app/user/controller"
	userRepo "github.com/BangNopall/paskihub-be/internal/app/user/repository"
	userSvc "github.com/BangNopall/paskihub-be/internal/app/user/service"

	eventCtr "github.com/BangNopall/paskihub-be/internal/app/event/controller"
	eventRepo "github.com/BangNopall/paskihub-be/internal/app/event/repository"
	eventSvc "github.com/BangNopall/paskihub-be/internal/app/event/service"

	walletCtr "github.com/BangNopall/paskihub-be/internal/app/wallet/controller"
	walletRepo "github.com/BangNopall/paskihub-be/internal/app/wallet/repository"
	walletSvc "github.com/BangNopall/paskihub-be/internal/app/wallet/service"
)

type Server interface {
	Start(port string)
	MountMiddlewares()
	RegistCustomValidation()
	MountRoutes(db *gorm.DB)
}

type httpServer struct {
	app       *fiber.App
	scheduler *cron.Cron
	validator *validator.Validate
}

func NewHttpServer() Server {
	app := fiber.New()
	scheduler := cron.New()
	validator := validator.New()

	return &httpServer{
		app:       app,
		scheduler: scheduler,
		validator: validator,
	}
}

func (s *httpServer) Start(port string) {
	if port[0] != ':' {
		port = ":" + port
	}

	s.app.Static("/public", "./public")
	err := s.app.Listen(port)

	if err != nil {
		log.Fatal(log.LogInfo{
			"error": err.Error(),
		}, "[SERVER][Start] failed to start server")
	}
}

func (s *httpServer) MountMiddlewares() {
	s.app.Use(middlewares.CORS())
	s.app.Use(middlewares.ApiKey())
}

func (s *httpServer) RegistCustomValidation() {
	s.validator.RegisterValidation("alphnumsympace", validators.Alphnumsympace)
	s.validator.RegisterValidation("plusnumeric", validators.Plusnumeric)
	s.validator.RegisterValidation("validdate", validators.DateValidation)
	s.validator.RegisterValidation("validpassword", validators.PasswordValidation)
}

func (s *httpServer) MountRoutes(db *gorm.DB) {
	uuid := uuid.UUID
	bcrypt := bcrypt.Bcrypt
	gomail := gomail.Gomail
	jwt := jwt.Jwt
	timePkg := timePkg.Time
	redis := redis.NewRedisClient()

	// Repository
	userRepo := userRepo.NewUserRepository(db)
	eventRepo := eventRepo.NewEventRepository(db)
	walletRepo := walletRepo.NewWalletRepository(db)

	// middleware
	middleware := middlewares.NewMiddleware(
		jwt,
		userRepo,
		redis,
	)

	// Service
	userSvc := userSvc.NewUserService(userRepo, uuid, bcrypt, timePkg, gomail, jwt, redis, time.Second*15)
	eventSvc := eventSvc.NewEventService(eventRepo, walletRepo, uuid, timePkg, time.Second*15)
	walletSvc := walletSvc.NewWalletService(walletRepo, eventRepo, uuid, time.Second*15)

	// Controller
	userCtr.InitUserController(userSvc, s.app, middleware, redis)
	eventCtr.InitEventController(eventSvc, s.app, middleware, redis)
	walletCtr.InitWalletController(walletSvc, s.app, middleware, redis)

	// cronjob
	_, err := s.scheduler.AddFunc("0 0 * * 0", userSvc.DeleteUnverifiedUsers)

	if err != nil {
		log.Fatal(log.LogInfo{
			"error": err.Error(),
		}, "[HTTP SERVER][Mount routes] failed to add cron job delete unverified users")
	}

	s.scheduler.Start()
}
