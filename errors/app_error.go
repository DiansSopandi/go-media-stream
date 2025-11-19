package errors

import (
	"fmt"
	"net/http"

	"github.com/DiansSopandi/media_stream/pkg"
)

type DetailResponse struct {
	Path       string `json:"path" example:"/api/v1/path"`
	Param      string `json:"param" example:"?page=1&limit=10"`
	StatusCode int    `json:"status_code" example:"200"`
	Method     string `json:"method" example:"GET"`
	Status     string `json:"status" example:"success_ok"`
} // @name	DetailResponse

type AppErrorResponse struct {
	Details    DetailResponse `json:"details"`
	Success    bool           `json:"success" example:"true"`
	Data       any            `json:"data" swaggertype:"array,object"`
	Errors     any            `json:"errors" swaggertype:"array,object"`
	Message    interface{}    `json:"message" example:"API Message"`
	LogMessage string         `json:"log_message" example:"Log message for internal user"`
} // @name ResponseApi

func (e *AppErrorResponse) Error() string {
	if msg, ok := e.Message.(string); ok {
		return msg
	}
	return fmt.Sprintf("%v", e.Message)
	// return e.Message.(string)
}

func NewAppErrorResponse(code string, statusCode int, logMessage string, status string) *AppErrorResponse {
	message, ok := ErrorCodeMap[code]

	if !ok {
		message = "Unknown error"
	}

	return &AppErrorResponse{
		Details: DetailResponse{
			Path:       "",
			Param:      "",
			StatusCode: statusCode,
			Method:     "",
			Status:     status,
		},
		Success:    false,
		Data:       nil,
		Errors:     nil,
		Message:    message,
		LogMessage: logMessage,
	}
}

func UserNotFound(logMessage string) *AppErrorResponse {
	return NewAppErrorResponse("USER_NOT_FOUND", http.StatusNotFound, logMessage, string(pkg.ApiStatusErrorNotFound))
}

func InvalidCredential(logMessage string) *AppErrorResponse {
	return NewAppErrorResponse("INVALID_CREDENTIAL", http.StatusUnauthorized, logMessage, string(pkg.ApiStatusErrorUnauthorized))
}

func DatabaseError(logMessage string) *AppErrorResponse {
	return NewAppErrorResponse("DB_ERROR", http.StatusInternalServerError, logMessage, string(pkg.ApiStatusErrorInternalServerError))
}

func InternalError(logMessage string) *AppErrorResponse {
	return NewAppErrorResponse("INTERNAL_ERROR", http.StatusInternalServerError, logMessage, string(pkg.ApiStatusErrorInternalServerError))
}

func RoleNotFound(logMessage string) *AppErrorResponse {
	return NewAppErrorResponse("ROLE_NOT_FOUND", http.StatusNotFound, logMessage, string(pkg.ApiStatusErrorNotFound))
}

func RoleValidationFailed(logMessage string) *AppErrorResponse {
	return NewAppErrorResponse("ROLE_VALDATION_FAILED", http.StatusInternalServerError, logMessage, string(pkg.ApiStatusErrorInternalServerError))
}

func PermissionDenied(logMessage string) *AppErrorResponse {
	return NewAppErrorResponse("PERMISSION_DENIED", http.StatusForbidden, logMessage, string(pkg.ApiStatusErrorForbidden))
}

func InvalidInput(logMessage string) *AppErrorResponse {
	return NewAppErrorResponse("INVALID_INPUT", http.StatusBadRequest, logMessage, string(pkg.ApiStatusErrorBadRequest))
}

func EmailAlreadyExists(logMessage string) *AppErrorResponse {
	return NewAppErrorResponse("EMAIL_ALREADY_EXISTS", http.StatusConflict, logMessage, string(pkg.ApiStatusErrorConflict))
}

func UsernameAlreadyExists(logMessage string) *AppErrorResponse {
	return NewAppErrorResponse("USERNAME_ALREADY_EXISTS", http.StatusConflict, logMessage, string(pkg.ApiStatusErrorConflict))
}

func PhoneAlreadyExists(logMessage string) *AppErrorResponse {
	return NewAppErrorResponse("PHONE_ALREADY_EXISTS", http.StatusConflict, logMessage, string(pkg.ApiStatusErrorConflict))
}

func Unauthorized(logMessage string) *AppErrorResponse {
	return NewAppErrorResponse("UNAUTHORIZED", http.StatusUnauthorized, logMessage, string(pkg.ApiStatusErrorUnauthorized))
}

func ResourceNotFound(logMessage string) *AppErrorResponse {
	return NewAppErrorResponse("RESOURCE_NOT_FOUND", http.StatusNotFound, logMessage, string(pkg.ApiStatusErrorNotFound))
}

func InvalidToken(logMessage string) *AppErrorResponse {
	return NewAppErrorResponse("INVALID_TOKEN", http.StatusUnauthorized, logMessage, string(pkg.ApiStatusErrorUnauthorized))
}

func PasswordMismatch(logMessage string) *AppErrorResponse {
	return NewAppErrorResponse("PASSWORD_MISMATCH", http.StatusBadRequest, logMessage, string(pkg.ApiStatusErrorBadRequest))
}

func UserAlreadyExists(logMessage string) *AppErrorResponse {
	return NewAppErrorResponse("USER_ALREADY_EXISTS", http.StatusConflict, logMessage, string(pkg.ApiStatusErrorConflict))
}

func RoleAlreadyExists(logMessage string) *AppErrorResponse {
	return NewAppErrorResponse("ROLE_ALREADY_EXISTS", http.StatusConflict, logMessage, string(pkg.ApiStatusErrorConflict))
}

func EmailNotVerified(logMessage string) *AppErrorResponse {
	return NewAppErrorResponse("EMAIL_NOT_VERIFIED", http.StatusForbidden, logMessage, string(pkg.ApiStatusErrorForbidden))
}

func PhoneNotVerified(logMessage string) *AppErrorResponse {
	return NewAppErrorResponse("PHONE_NOT_VERIFIED", http.StatusForbidden, logMessage, string(pkg.ApiStatusErrorForbidden))
}

func AccountLocked(logMessage string) *AppErrorResponse {
	return NewAppErrorResponse("ACCOUNT_LOCKED", http.StatusForbidden, logMessage, string(pkg.ApiStatusErrorForbidden))
}

func TooManyRequests(logMessage string) *AppErrorResponse {
	return NewAppErrorResponse("TOO_MANY_REQUESTS", http.StatusTooManyRequests, logMessage, string(pkg.ApiStatusErrorTooManyRequests))
}

func InvalidFileType(logMessage string) *AppErrorResponse {
	return NewAppErrorResponse("INVALID_FILE_TYPE", http.StatusBadRequest, logMessage, string(pkg.ApiStatusErrorBadRequest))
}

func FileTooLarge(logMessage string) *AppErrorResponse {
	return NewAppErrorResponse("FILE_TOO_LARGE", http.StatusBadRequest, logMessage, string(pkg.ApiStatusErrorBadRequest))
}

func OperationNotAllowed(logMessage string) *AppErrorResponse {
	return NewAppErrorResponse("OPERATION_NOT_ALLOWED", http.StatusMethodNotAllowed, logMessage, string(pkg.ApiStatusErrorMethodNotAllowed))
}

func ResourceConflict(logMessage string) *AppErrorResponse {
	return NewAppErrorResponse("RESOURCE_CONFLICT", http.StatusConflict, logMessage, string(pkg.ApiStatusErrorConflict))
}
func InvalidRequest(logMessage string) *AppErrorResponse {
	return NewAppErrorResponse("INVALID_REQUEST", http.StatusBadRequest, logMessage, string(pkg.ApiStatusErrorBadRequest))
}
func InvalidRequestWithMessage(message string) *AppErrorResponse {
	return &AppErrorResponse{
		Details: DetailResponse{
			Path:       "",
			Param:      "",
			StatusCode: http.StatusBadRequest,
			Method:     "",
			Status:     "",
		},
		Success: false,
		Data:    nil,
		Errors:  nil,
		Message: message,
	}
}

func InvalidRequestWithCode(code string) *AppErrorResponse {
	return &AppErrorResponse{
		Details: DetailResponse{
			Path:       "",
			Param:      "",
			StatusCode: http.StatusBadRequest,
			Method:     "",
			Status:     "",
		},
		Success: false,
		Data:    nil,
		Errors:  nil,
		Message: "Invalid request",
	}
}

func InvalidRequestWithCodeAndMessage(code, message string) *AppErrorResponse {
	return &AppErrorResponse{
		Details: DetailResponse{
			Path:       "",
			Param:      "",
			StatusCode: http.StatusBadRequest,
			Method:     "",
			Status:     "",
		},
		Success: false,
		Data:    nil,
		Errors:  nil,
		Message: message,
	}
}

func InvalidRequestWithCodeMessageAndStatus(code, message string, status int) *AppErrorResponse {
	return &AppErrorResponse{
		Details: DetailResponse{
			Path:       "",
			Param:      "",
			StatusCode: status,
			Method:     "",
			Status:     "",
		},
		Success: false,
		Data:    nil,
		Errors:  nil,
		Message: message,
	}
}
