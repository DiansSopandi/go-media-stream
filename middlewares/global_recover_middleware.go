package middlewares

import (
	"fmt"
	"runtime/debug"

	"github.com/DiansSopandi/media_stream/errors"
	"github.com/DiansSopandi/media_stream/pkg"
	"github.com/gofiber/fiber/v2"
)

// GlobalRecoveryMiddleware menangkap panic di seluruh route
func GlobalRecoveryMiddleware(c *fiber.Ctx) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// Log error + stack trace
			// log.Printf("[PANIC RECOVERED] %v\n%s", r, debug.Stack())
			stackTrace := string(debug.Stack())

			logMessage := fmt.Sprintf("[PANIC RECOVERED] %v\n%s", r, stackTrace)
			pkg.CreateAccessLog(c, "[ACCESS:API][PANIC]", fiber.StatusInternalServerError, logMessage)

			err = &errors.AppErrorResponse{
				Details: errors.DetailResponse{
					Path:       c.Request().URI().String(),
					Method:     string(c.Request().Header.Method()),
					StatusCode: fiber.StatusInternalServerError,
					Status:     string(pkg.ApiStatusErrorInternalServerError),
				},
				Success:    false,
				Data:       nil,
				Message:    fmt.Sprintf("%v", r), // Message untuk client
				LogMessage: logMessage,           // Message untuk internal log
			}
		}
	}()
	return c.Next()
}
