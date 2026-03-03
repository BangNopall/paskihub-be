package controller

import (
	"net/http"
	"strconv"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/internal/middlewares"
	"github.com/BangNopall/paskihub-be/pkg/helpers/http/response"
	"github.com/BangNopall/paskihub-be/pkg/redis"
	"github.com/gofiber/fiber/v2"
)

type walletController struct {
	walletSvc contracts.WalletService
	redis     redis.RedisInterface
}

func InitWalletController(
	walletSvc contracts.WalletService,
	router fiber.Router,
	middleware *middlewares.Middleware,
	redis redis.RedisInterface,
) {
	walletController := &walletController{
		walletSvc: walletSvc,
		redis:     redis,
	}

	walletRouter := router.Group("/api/v1/wallets")

	// Organizer Routes
	walletRouter.Get("/:eventId", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, walletController.GetWalletInfo)
	walletRouter.Get("/:eventId/logs", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, walletController.GetTransactionLogs)
	walletRouter.Post("/:eventId/topup", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, walletController.RequestTopUp)

	// Admin Route
	walletRouter.Put("/transactions/:transactionId/approve", middleware.Authentication, middleware.RateLimiter(), middleware.AuthAdmin, walletController.ApproveTopUp)
	walletRouter.Put("/transactions/:transactionId/reject", middleware.Authentication, middleware.RateLimiter(), middleware.AuthAdmin, walletController.RejectTopUp)
}

func (c *walletController) GetWalletInfo(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to get wallet info"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	eventId := ctx.Params("eventId")
	userId := ctx.Locals("id").(string)

	res, err = c.walletSvc.GetWalletInfo(ctx.Context(), eventId, userId)
	code = domain.GetCode(err)
	if err != nil {
		return nil
	}

	message = "success to get wallet info"
	return nil
}

func (c *walletController) GetTransactionLogs(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to get transaction logs"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	eventId := ctx.Params("eventId")
	userId := ctx.Locals("id").(string)

	res, err = c.walletSvc.GetTransactionLogs(ctx.Context(), eventId, userId)
	code = domain.GetCode(err)
	if err != nil {
		return nil
	}

	message = "success to get transaction logs"
	return nil
}

func (c *walletController) RequestTopUp(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to request top up"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	eventId := ctx.Params("eventId")
	userId := ctx.Locals("id").(string)

	amountStr := ctx.FormValue("amount")
	couponCode := ctx.FormValue("coupon_code")

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		code = http.StatusBadRequest
		err = domain.ErrBadRequest
		return nil
	}

	file, err := ctx.FormFile("proof")
	if err != nil {
		code = http.StatusBadRequest
		err = domain.ErrBadRequest
		return nil
	}

	req := &dto.TopUpRequest{
		Amount:     amount,
		CouponCode: couponCode,
	}

	err = c.walletSvc.RequestTopUp(ctx.Context(), eventId, userId, req, file)
	code = domain.GetCode(err)
	if err != nil {
		return nil
	}

	message = "success to request top up"
	return nil
}

func (c *walletController) ApproveTopUp(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to approve top up request"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	transactionId := ctx.Params("transactionId")
	adminUserIdStr := ctx.Locals("id").(string)

	err = c.walletSvc.ApproveTopUp(ctx.Context(), transactionId, adminUserIdStr)
	code = domain.GetCode(err)
	if err != nil {
		return nil
	}

	message = "success to approve top up request"
	return nil
}

func (c *walletController) RejectTopUp(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusBadRequest
		res     interface{}
		message string = "failed to reject top up request"
	)

	sendResp := func() {
		response.Send(ctx, code, message, res, err)
	}
	defer sendResp()

	transactionId := ctx.Params("transactionId")
	adminUserIdStr := ctx.Locals("id").(string)

	err = c.walletSvc.RejectTopUp(ctx.Context(), transactionId, adminUserIdStr)
	code = domain.GetCode(err)
	if err != nil {
		return nil
	}

	message = "success to reject top up request"
	return nil
}
