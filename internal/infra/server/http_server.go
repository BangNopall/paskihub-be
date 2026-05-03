package server

import (
	"time"

	"github.com/BangNopall/paskihub-be/internal/infra/env"
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
	"github.com/gofiber/swagger"
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

	assessmentCtr "github.com/BangNopall/paskihub-be/internal/app/assessment/controller"
	assessmentRepo "github.com/BangNopall/paskihub-be/internal/app/assessment/repository"
	assessmentSvc "github.com/BangNopall/paskihub-be/internal/app/assessment/service"

	pProfileCtr "github.com/BangNopall/paskihub-be/internal/app/participant_profile/controller"
	pProfileRepo "github.com/BangNopall/paskihub-be/internal/app/participant_profile/repository"
	pProfileSvc "github.com/BangNopall/paskihub-be/internal/app/participant_profile/service"

	pTeamCtr "github.com/BangNopall/paskihub-be/internal/app/participant_team/controller"
	pTeamRepo "github.com/BangNopall/paskihub-be/internal/app/participant_team/repository"
	pTeamSvc "github.com/BangNopall/paskihub-be/internal/app/participant_team/service"

	pEventCtr "github.com/BangNopall/paskihub-be/internal/app/participant_event/controller"
	pEventRepo "github.com/BangNopall/paskihub-be/internal/app/participant_event/repository"
	pEventSvc "github.com/BangNopall/paskihub-be/internal/app/participant_event/service"

	pAssessCtr "github.com/BangNopall/paskihub-be/internal/app/participant_assessment/controller"
	pAssessRepo "github.com/BangNopall/paskihub-be/internal/app/participant_assessment/repository"
	pAssessSvc "github.com/BangNopall/paskihub-be/internal/app/participant_assessment/service"

	eoTeamCtr "github.com/BangNopall/paskihub-be/internal/app/eo_team/controller"
	eoTeamRepo "github.com/BangNopall/paskihub-be/internal/app/eo_team/repository"
	eoTeamSvc "github.com/BangNopall/paskihub-be/internal/app/eo_team/service"

	dashboardCtr "github.com/BangNopall/paskihub-be/internal/app/dashboard/controller"
	dashboardRepo "github.com/BangNopall/paskihub-be/internal/app/dashboard/repository"
	dashboardSvc "github.com/BangNopall/paskihub-be/internal/app/dashboard/service"

	eoCtr "github.com/BangNopall/paskihub-be/internal/app/eo/controller"
	eoRepo "github.com/BangNopall/paskihub-be/internal/app/eo/repository"
	eoSvc "github.com/BangNopall/paskihub-be/internal/app/eo/service"

	systemSettingCtr "github.com/BangNopall/paskihub-be/internal/app/system_setting/controller"
	systemSettingRepo "github.com/BangNopall/paskihub-be/internal/app/system_setting/repository"
	systemSettingSvc "github.com/BangNopall/paskihub-be/internal/app/system_setting/service"
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
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // 10MB limit for file uploads
	})
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

	err := s.app.Listen(port)

	if err != nil {
		log.Fatal(log.LogInfo{
			"error": err.Error(),
		}, "[SERVER][Start] failed to start server")
	}
}

func (s *httpServer) MountMiddlewares() {
	if env.AppEnv.AppEnv == "development" {
		s.app.Get("/swagger/*", swagger.HandlerDefault)
	}

	s.app.Use(middlewares.CORS())
	s.app.Static("/public", "./public")
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
	formPenilaianRepo := assessmentRepo.NewFormPenilaianRepository(db)
	rekapRepo := assessmentRepo.NewRekapRepository(db)
	juryRepo := assessmentRepo.NewJuryRepository(db)

	assessmentRepoIns := assessmentRepo.NewAssessmentRepository(db)

	pProfileRepoIns := pProfileRepo.NewParticipantProfileRepository(db)
	pTeamRepoIns := pTeamRepo.NewParticipantTeamRepository(db)
	pEventRepoIns := pEventRepo.NewParticipantEventRepository(db)
	pAssessRepoIns := pAssessRepo.NewParticipantAssessmentRepository(db)
	eoTeamRepoIns := eoTeamRepo.NewEOTeamRepository(db)
	dashboardRepoIns := dashboardRepo.NewDashboardRepository(db)
	eoRepoIns := eoRepo.NewEORepository(db)
	settingRepoIns := systemSettingRepo.NewSystemSettingRepository(db)

	// middleware
	middleware := middlewares.NewMiddleware(
		jwt,
		userRepo,
		redis,
	)

	// Service
	userSvc := userSvc.NewUserService(userRepo, uuid, bcrypt, timePkg, gomail, jwt, redis, time.Second*15)
	eventSvc := eventSvc.NewEventService(eventRepo, walletRepo, uuid, timePkg, time.Second*15)
	walletSvc := walletSvc.NewWalletService(walletRepo, eventRepo, settingRepoIns, uuid, time.Second*15)
	formPenilaianSvc := assessmentSvc.NewFormPenilaianService(formPenilaianRepo, db, s.validator)
	rekapSvc := assessmentSvc.NewRekapService(rekapRepo)
	jurySvc := assessmentSvc.NewJuryService(juryRepo)
	assessmentSvcIns := assessmentSvc.NewAssessmentService(assessmentRepoIns)

	pProfileSvcIns := pProfileSvc.NewParticipantProfileService(pProfileRepoIns)
	pTeamSvcIns := pTeamSvc.NewParticipantTeamService(pTeamRepoIns, pProfileRepoIns)
	pEventSvcIns := pEventSvc.NewParticipantEventService(pEventRepoIns)
	pAssessSvcIns := pAssessSvc.NewParticipantAssessmentService(pAssessRepoIns)
	eoTeamSvcIns := eoTeamSvc.NewEOTeamService(eoTeamRepoIns, walletRepo, settingRepoIns, db)
	dashboardSvcIns := dashboardSvc.NewDashboardService(dashboardRepoIns)
	eoSvcIns := eoSvc.NewEOService(eoRepoIns, userRepo, bcrypt, uuid, time.Second*15)
	settingSvcIns := systemSettingSvc.NewSystemSettingService(settingRepoIns, time.Second*15)

	// Controller
	userCtr.InitUserController(userSvc, s.app, middleware, redis)
	eventCtr.InitEventController(eventSvc, s.app, middleware, redis)
	walletCtr.InitWalletController(walletSvc, s.app, middleware, redis)
	assessmentCtr.InitFormPenilaianController(formPenilaianSvc, s.app, middleware)
	assessmentCtr.InitRekapController(rekapSvc, s.app, middleware)
	assessmentCtr.InitJuryController(jurySvc, s.app, middleware)
	assessmentCtr.InitAssessmentController(assessmentSvcIns, s.app, middleware)
	eoTeamCtr.InitEOTeamController(eoTeamSvcIns, s.app, middleware)
	dashboardCtr.InitDashboardController(dashboardSvcIns, s.app, middleware)
	systemSettingCtr.InitSystemSettingController(settingSvcIns, s.app, middleware)

	eoControllerIns := eoCtr.NewEOController(eoSvcIns, s.validator)
	eoControllerIns.Route(s.app, middleware)

	pesertaGrp := s.app.Group("/api/v1/peserta", middleware.Authentication, middleware.RateLimiter(), middleware.AuthPeserta)

	pProfileController := pProfileCtr.NewParticipantProfileController(pProfileSvcIns, s.validator)
	pProfileController.Route(pesertaGrp)

	pTeamController := pTeamCtr.NewParticipantTeamController(pTeamSvcIns, s.validator)
	pTeamController.Route(pesertaGrp)

	pEventController := pEventCtr.NewParticipantEventController(pEventSvcIns, s.validator)
	pEventController.Route(pesertaGrp)

	pAssessController := pAssessCtr.NewParticipantAssessmentController(pAssessSvcIns)
	pAssessController.Route(pesertaGrp)

	// cronjob
	_, err := s.scheduler.AddFunc("0 0 * * 0", userSvc.DeleteUnverifiedUsers)

	if err != nil {
		log.Fatal(log.LogInfo{
			"error": err.Error(),
		}, "[HTTP SERVER][Mount routes] failed to add cron job delete unverified users")
	}

	s.scheduler.Start()
}
