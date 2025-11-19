package middlewares

import (
	"github.com/DiansSopandi/media_stream/errors"
	"github.com/DiansSopandi/media_stream/pkg"
	"github.com/gofiber/fiber/v2"
)

type DetailResponse struct {
	Path       string `json:"path" example:"/api/v1/path"`
	Param      string `json:"param" example:"?page=1&limit=10"`
	StatusCode int    `json:"status_code" example:"200"`
	Method     string `json:"method" example:"GET"`
	Status     string `json:"status" example:"success_ok"`
} // @name	DetailResponse

type ErrorResponse struct {
	Details    DetailResponse `json:"details"`
	Success    bool           `json:"success" example:"true"`
	Data       any            `json:"data" swaggertype:"array,object"`
	Errors     any            `json:"errors" swaggertype:"array,object"`
	Message    interface{}    `json:"message" example:"API Message"`
	LogMessage string         `json:"-" example:"Log message for internal user"`
} // @name ResponseApi

func ErrorHandler(c *fiber.Ctx, err error) error {
	detail := DetailResponse{
		StatusCode: int(pkg.HttpStatusInternalServerError),
		Path:       c.Request().URI().String(),
		Method:     string(c.Request().Header.Method()),
		Status:     string(pkg.ApiStatusErrorInternalServerError),
	}

	res := ErrorResponse{
		Details:    detail,
		Success:    false,
		Data:       nil,
		Message:    err.Error(), // errors.ErrorCodeMap["INTERNAL_ERROR"],
		LogMessage: err.Error(),
	}

	if appErr, ok := err.(*errors.AppErrorResponse); ok {
		res.Details.StatusCode = appErr.Details.StatusCode
		res.Details.Status = appErr.Details.Status
		res.Message = appErr.Message
		res.LogMessage = appErr.LogMessage

		pkg.CreateAccessLog(c, "[ACCESS:API][ERROR]", res.Details.StatusCode, appErr.LogMessage)
	} else if e, ok := err.(*fiber.Error); ok {
		res.Details.StatusCode = e.Code
		res.Message = e.Message
		res.LogMessage = e.Error()
	}

	if _, ok := err.(*errors.AppErrorResponse); !ok {
		pkg.CreateAccessLog(c, "[ACCESS:API][ERROR]", res.Details.StatusCode, err.Error())
	}

	return c.Status(res.Details.StatusCode).JSON(res)
}
