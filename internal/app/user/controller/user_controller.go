package controller

import (
	"net/http"
	"strings"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/internal/middlewares"
	"github.com/BangNopall/paskihub-be/pkg/helpers/http/response"
	"github.com/BangNopall/paskihub-be/pkg/redis"
	"github.com/gofiber/fiber/v2"
)

type userController struct {
	userSvc contracts.UserService
	redis   redis.RedisInterface
}

func InitUserController(
	userSvc contracts.UserService,
	router fiber.Router,
	middleware *middlewares.Middleware,
	redis redis.RedisInterface,
) {
	userController := &userController{
		userSvc: userSvc,
		redis:   redis,
	}

	userRouter := router.Group("/api/v1/users")
	userRouter.Post("/register/:role", middleware.RateLimiter(), userController.Register)
	userRouter.Post("/login", middleware.RateLimiter(), userController.Login)
	userRouter.Post("/logout", middleware.Authentication, userController.Logout)
	userRouter.Get("/verify-email/:email/:emailVerPass", middleware.RateLimiter(), userController.VerifyEmail)
	userRouter.Put("/reset-password/:token", middleware.RateLimiter(), userController.ResetPassword)
	userRouter.Post("/forgot-password", middleware.RateLimiter(), userController.ForgotPassword)
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with specific role (eo or peserta)
// @Tags Users
// @Accept json
// @Produce json
// @Param role path string true "User Role"
// @Param user body dto.UserRegister true "User Data"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/users/register/{role} [post]
func (c *userController) Register(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to register account"
	)

	sendResp := func() {
		response.Send(
			ctx,
			code,
			message,
			res,
			err,
		)
	}

	defer sendResp()

	var user dto.UserRegister

	role := ctx.Params("role")

	err = ctx.BodyParser(&user)

	if err != nil {
		return nil
	}

	err = c.userSvc.Register(ctx.Context(), role, user)
	code = domain.GetCode(err)

	if err != nil {
		return nil
	}

	message = "success to register account, please check your email"
	return nil
}

// Login godoc
// @Summary User login
// @Description Login with email and password
// @Tags Users
// @Accept json
// @Produce json
// @Param user body dto.UserLogin true "Login Credentials"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/users/login [post]
func (c *userController) Login(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to login with email"
	)

	sendResp := func() {
		response.Send(
			ctx,
			code,
			message,
			res,
			err,
		)
	}

	defer sendResp()

	var user dto.UserLogin

	err = ctx.BodyParser(&user)

	if err != nil {
		return nil
	}

	res, err = c.userSvc.Login(ctx.Context(), user)
	code = domain.GetCode(err)

	if err != nil {
		return nil
	}

	message = "success to login with email"
	return nil
}

// VerifyEmail godoc
// @Summary Verify user email
// @Description Verify email using token sent to user email
// @Tags Users
// @Produce json
// @Param email path string true "User Email"
// @Param emailVerPass path string true "Verification Password"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/users/verify-email/{email}/{emailVerPass} [get]
func (c *userController) VerifyEmail(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to verify email"
	)

	sendResp := func() {
		response.Send(
			ctx,
			code,
			message,
			res,
			err,
		)
	}

	defer sendResp()

	email := ctx.Params("email")
	emailVerPass := ctx.Params("emailVerPass")

	err = c.userSvc.VerifyEmail(ctx.Context(), email, emailVerPass)
	code = domain.GetCode(err)

	if err != nil {
		return nil
	}

	message = "success to verify email"
	return nil
}

// ResetPassword godoc
// @Summary Reset user password
// @Description Reset password using reset token
// @Tags Users
// @Accept json
// @Produce json
// @Param token path string true "Reset Token"
// @Param user body dto.UserResetPassword true "New Password"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/users/reset-password/{token} [put]
func (c *userController) ResetPassword(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to reset password"
	)

	sendResp := func() {
		response.Send(
			ctx,
			code,
			message,
			res,
			err,
		)
	}

	defer sendResp()

	token := ctx.Params("token")
	var user dto.UserResetPassword

	err = ctx.BodyParser(&user)

	if err != nil {
		return nil
	}

	err = c.userSvc.ResetPassword(ctx.Context(), user, token)
	code = domain.GetCode(err)

	if err != nil {
		return nil
	}

	message = "success to reset password"
	return nil
}

// ForgotPassword godoc
// @Summary Forgot user password
// @Description Request a password reset link to user email
// @Tags Users
// @Accept json
// @Produce json
// @Param user body dto.UserForgotPassword true "User Email"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/users/forgot-password [post]
func (c *userController) ForgotPassword(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to forgot password"
	)

	sendResp := func() {
		response.Send(
			ctx,
			code,
			message,
			res,
			err,
		)
	}

	defer sendResp()

	var user dto.UserForgotPassword

	err = ctx.BodyParser(&user)

	if err != nil {
		return nil
	}

	err = c.userSvc.ForgotPassword(ctx.Context(), user)
	code = domain.GetCode(err)

	if err != nil {
		return nil
	}

	message = "success to forgot password"
	return nil
}

// Logout godoc
// @Summary User logout
// @Description Logout user and invalidate token
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/users/logout [post]
func (c *userController) Logout(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to logout"
	)

	sendResp := func() {
		response.Send(
			ctx,
			code,
			message,
			res,
			err,
		)
	}

	defer sendResp()

	bearerToken := ctx.Get("Authorization")

	token := strings.Split(bearerToken, " ")[1]

	err = c.userSvc.Logout(ctx.Context(), token)
	code = domain.GetCode(err)

	if err != nil {
		return nil
	}

	message = "success to logout"
	return nil
}
