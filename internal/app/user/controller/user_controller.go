package controller

import (
	"net/http"
	"strings"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/BangNopall/paskihub-be/internal/middlewares"
	"github.com/BangNopall/paskihub-be/pkg/helpers/http/response"
	"github.com/BangNopall/paskihub-be/pkg/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type userController struct {
	userSvc  contracts.UserService
	eventSvc contracts.EventService
	redis    redis.RedisInterface
}

func InitUserController(
	userSvc contracts.UserService,
	eventSvc contracts.EventService,
	router fiber.Router,
	middleware *middlewares.Middleware,
	redis redis.RedisInterface,
) {
	userController := &userController{
		userSvc:  userSvc,
		eventSvc: eventSvc,
		redis:    redis,
	}

	userRouter := router.Group("/api/v1/users")
	userRouter.Post("/register/:role", middleware.RateLimiter(), userController.Register)
	userRouter.Post("/login", middleware.RateLimiter(), userController.Login)
	userRouter.Post("/logout", middleware.Authentication, userController.Logout)
	userRouter.Get("/verify-email/:email/:emailVerPass", middleware.RateLimiter(), userController.VerifyEmail)
	userRouter.Put("/reset-password/:token", middleware.RateLimiter(), userController.ResetPassword)
	userRouter.Post("/forgot-password", middleware.RateLimiter(), userController.ForgotPassword)
	userRouter.Get("/show/:userId", middleware.Authentication, middleware.RateLimiter(), userController.GetUserById)

	// Admin Routes
	adminRouter := router.Group("/api/v1/admin")
	adminRouter.Use(middleware.Authentication, middleware.AuthAdmin)
	adminRouter.Get("/users", userController.GetUsers)
	adminRouter.Get("/admins", userController.GetAdmins)
	adminRouter.Post("/admins", userController.CreateAdmin)
	adminRouter.Delete("/admins/:id", userController.DeleteAdmin)
	adminRouter.Post("/admins/:id/reset-password", userController.ResetAdminPassword)
	adminRouter.Get("/users/:userId", userController.GetUserDetailAdmin)
	adminRouter.Put("/users/:userId/ban", userController.BanUser)
	adminRouter.Put("/users/:userId/unban", userController.UnbanUser)
	adminRouter.Put("/users/:userId/verify", userController.VerifyUser)
	adminRouter.Put("/users/:userId/archive", userController.ArchiveOrganizer)
	adminRouter.Put("/users/:userId/unarchive", userController.UnarchiveOrganizer)
	adminRouter.Put("/events/:eventId/status", userController.UpdateEventStatus)
}

// CreateAdmin godoc
// @Summary Create a new admin
// @Description Create a new admin staff account (Admin only)
// @Tags Admin
// @Security ApiKeyAuth && BearerAuth
// @Accept json
// @Produce json
// @Param admin body dto.AdminCreateRequest true "Admin Data"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/admins [post]
func (c *userController) CreateAdmin(ctx *fiber.Ctx) error {
	var req dto.AdminCreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		response.Send(ctx, http.StatusBadRequest, "invalid request", nil, err)
		return nil
	}

	err := c.userSvc.CreateAdmin(ctx.Context(), req)
	if err != nil {
		response.Send(ctx, domain.GetCode(err), "failed to create admin", nil, err)
		return nil
	}

	response.Send(ctx, http.StatusOK, "success to create admin", nil, nil)
	return nil
}

// DeleteAdmin godoc
// @Summary Delete admin account
// @Description Permanently remove an admin's access (Admin only)
// @Tags Admin
// @Security ApiKeyAuth && BearerAuth
// @Produce json
// @Param id path string true "Admin User ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/admins/{id} [delete]
func (c *userController) DeleteAdmin(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	userId, err := uuid.Parse(idStr)
	if err != nil {
		response.Send(ctx, http.StatusBadRequest, "invalid user id", nil, err)
		return nil
	}

	err = c.userSvc.DeleteAdmin(ctx.Context(), userId)
	if err != nil {
		response.Send(ctx, domain.GetCode(err), "failed to delete admin", nil, err)
		return nil
	}

	response.Send(ctx, http.StatusOK, "success to delete admin", nil, nil)
	return nil
}

// ResetAdminPassword godoc
// @Summary Reset admin password
// @Description Manually trigger a password reset for an admin staff (Admin only)
// @Tags Admin
// @Security ApiKeyAuth && BearerAuth
// @Produce json
// @Param id path string true "Admin User ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/admins/{id}/reset-password [post]
func (c *userController) ResetAdminPassword(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	userId, err := uuid.Parse(idStr)
	if err != nil {
		response.Send(ctx, http.StatusBadRequest, "invalid user id", nil, err)
		return nil
	}

	err = c.userSvc.ResetAdminPassword(ctx.Context(), userId)
	if err != nil {
		response.Send(ctx, domain.GetCode(err), "failed to reset admin password", nil, err)
		return nil
	}

	response.Send(ctx, http.StatusOK, "success to trigger password reset", nil, nil)
	return nil
}

// GetUserById godoc
// @Summary Get user by ID
// @Description Get user details by ID (Self access only)
// @Tags Users
// @Security ApiKeyAuth && BearerAuth
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/users/show/{userId} [get]
func (c *userController) GetUserById(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to fetch user"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	userIdStr := ctx.Params("userId")
	authUserIdStr := ctx.Locals("id").(string)

	if authUserIdStr != userIdStr {
		code = http.StatusForbidden
		message = "forbidden access"
		err = domain.ErrForbidden
		return nil
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return nil
	}

	res, err = c.userSvc.FetchUserById(ctx.Context(), userId)
	code = domain.GetCode(err)

	if err != nil {
		return nil
	}

	message = "success to fetch user"
	return nil
}

// GetUsers godoc
// @Summary Get all users
// @Description Retrieve a list of all users (Admin only)
// @Tags Admin
// @Security ApiKeyAuth && BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/users [get]
func (c *userController) GetUsers(ctx *fiber.Ctx) error {
	res, err := c.userSvc.FetchAllUsers(ctx.Context(), nil)
	if err != nil {
		response.Send(ctx, http.StatusInternalServerError, "failed to fetch users", nil, err)
		return nil
	}
	response.Send(ctx, http.StatusOK, "success to fetch users", res, nil)
	return nil
}

// GetAdmins godoc
// @Summary Get all admins
// @Description Retrieve a list of all admins (Admin only)
// @Tags Admin
// @Security ApiKeyAuth && BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/admins [get]
func (c *userController) GetAdmins(ctx *fiber.Ctx) error {
	role := string(enums.Admin)
	res, err := c.userSvc.FetchAllUsers(ctx.Context(), &role)
	if err != nil {
		response.Send(ctx, http.StatusInternalServerError, "failed to fetch admins", nil, err)
		return nil
	}
	response.Send(ctx, http.StatusOK, "success to fetch admins", res, nil)
	return nil
}

// GetUserDetailAdmin godoc
// @Summary Get admin user detail
// @Description Retrieve detailed info of a user (Admin only)
// @Tags Admin
// @Security ApiKeyAuth && BearerAuth
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/users/{userId} [get]
func (c *userController) GetUserDetailAdmin(ctx *fiber.Ctx) error {
	userIdStr := ctx.Params("userId")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		response.Send(ctx, http.StatusBadRequest, "invalid user id", nil, err)
		return nil
	}

	res, err := c.userSvc.FetchUserDetailAdmin(ctx.Context(), userId)
	if err != nil {
		response.Send(ctx, domain.GetCode(err), "failed to fetch user detail", nil, err)
		return nil
	}

	response.Send(ctx, http.StatusOK, "success to fetch user detail", res, nil)
	return nil
}

// BanUser godoc
// @Summary Ban user
// @Description Set user status to banned (Admin only)
// @Tags Admin
// @Security ApiKeyAuth && BearerAuth
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/users/{userId}/ban [put]
func (c *userController) BanUser(ctx *fiber.Ctx) error {
	userIdStr := ctx.Params("userId")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		response.Send(ctx, http.StatusBadRequest, "invalid user id", nil, err)
		return nil
	}

	err = c.userSvc.BanUser(ctx.Context(), userId)
	if err != nil {
		response.Send(ctx, domain.GetCode(err), "failed to ban user", nil, err)
		return nil
	}

	response.Send(ctx, http.StatusOK, "success to ban user", nil, nil)
	return nil
}

// UnbanUser godoc
// @Summary Unban user
// @Description Set user status to active (Admin only)
// @Tags Admin
// @Security ApiKeyAuth && BearerAuth
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/users/{userId}/unban [put]
func (c *userController) UnbanUser(ctx *fiber.Ctx) error {
	userIdStr := ctx.Params("userId")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		response.Send(ctx, http.StatusBadRequest, "invalid user id", nil, err)
		return nil
	}

	err = c.userSvc.UnbanUser(ctx.Context(), userId)
	if err != nil {
		response.Send(ctx, domain.GetCode(err), "failed to unban user", nil, err)
		return nil
	}

	response.Send(ctx, http.StatusOK, "success to unban user", nil, nil)
	return nil
}

// VerifyUser godoc
// @Summary Verify user email
// @Description Manually verify a user's email (Admin only)
// @Tags Admin
// @Security ApiKeyAuth && BearerAuth
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/users/{userId}/verify [put]
func (c *userController) VerifyUser(ctx *fiber.Ctx) error {
	userIdStr := ctx.Params("userId")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		response.Send(ctx, http.StatusBadRequest, "invalid user id", nil, err)
		return nil
	}

	err = c.userSvc.VerifyUser(ctx.Context(), userId)
	if err != nil {
		response.Send(ctx, domain.GetCode(err), "failed to verify user", nil, err)
		return nil
	}

	response.Send(ctx, http.StatusOK, "success to verify user", nil, nil)
	return nil
}

// ArchiveOrganizer godoc
// @Summary Archive organizer event
// @Description Set organizer's event status to archived (Admin only)
// @Tags Admin
// @Security ApiKeyAuth && BearerAuth
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/users/{userId}/archive [put]
func (c *userController) ArchiveOrganizer(ctx *fiber.Ctx) error {
	userIdStr := ctx.Params("userId")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		response.Send(ctx, http.StatusBadRequest, "invalid user id", nil, err)
		return nil
	}

	err = c.userSvc.ArchiveOrganizer(ctx.Context(), userId)
	if err != nil {
		response.Send(ctx, domain.GetCode(err), "failed to archive organizer", nil, err)
		return nil
	}

	response.Send(ctx, http.StatusOK, "success to archive organizer", nil, nil)
	return nil
}

// UnarchiveOrganizer godoc
// @Summary Unarchive organizer event
// @Description Set organizer's event status to open (Admin only)
// @Tags Admin
// @Security ApiKeyAuth && BearerAuth
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/users/{userId}/unarchive [put]
func (c *userController) UnarchiveOrganizer(ctx *fiber.Ctx) error {
	userIdStr := ctx.Params("userId")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		response.Send(ctx, http.StatusBadRequest, "invalid user id", nil, err)
		return nil
	}

	err = c.userSvc.UnarchiveOrganizer(ctx.Context(), userId)
	if err != nil {
		response.Send(ctx, domain.GetCode(err), "failed to unarchive organizer", nil, err)
		return nil
	}

	response.Send(ctx, http.StatusOK, "success to unarchive organizer", nil, nil)
	return nil
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
// @Security ApiKeyAuth
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
// @Security ApiKeyAuth
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
// @Security ApiKeyAuth
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
// @Security ApiKeyAuth
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
// @Security ApiKeyAuth
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
// @Security ApiKeyAuth && BearerAuth
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

