package response

import "github.com/gofiber/fiber/v2"

type Status string

const (
	// Success thrown when all went well
	Success Status = "success"

	/* Failed thrown when there is a problem from request data
	ex: (validation error, failed to bind json data, etc) */
	Fail Status = "fail"

	/* Error thrown when an error occured in processing the request aka server error */
	Error Status = "error"
)

type Response struct {
	Code    int         `json:"code"`
	Status  Status      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorDetail struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

func SendResp(ctx *fiber.Ctx, code int, status Status, message string, data interface{}) {
	ctx.Status(code).JSON(Response{
		Code:    code,
		Status:  status,
		Message: message,
		Data:    data,
	})
}

func SendErrResp(
	ctx *fiber.Ctx,
	code int,
	status Status,
	message string,
	err error,
) {
	errMsg := ""

	if err != nil {
		if code >= 500 {
			errMsg = "Internal Server Error"
		} else {
			errMsg = err.Error()
		}
	}

	ctx.Status(code).JSON(ErrorResponse{
		Error: ErrorDetail{
			Code:    string(status),
			Message: message,
			Details: errMsg,
		},
	})
}

func Send(
	ctx *fiber.Ctx,
	code int,
	message string,
	data interface{},
	err error,
) {
	status := GetStatus(code)

	if err != nil {
		SendErrResp(
			ctx,
			code,
			status,
			message,
			err,
		)
	} else {
		SendResp(
			ctx,
			code,
			status,
			message,
			data,
		)
	}

}

func GetStatus(code int) Status {
	if code >= 500 {
		return Error
	}

	if code >= 400 {
		return Fail
	}

	return Success
}
