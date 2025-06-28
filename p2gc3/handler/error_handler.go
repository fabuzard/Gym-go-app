// handler/error_handler.go
package handler

import (
	"fmt"
	"net/http"
	"p2gc3/dto"
	"p2gc3/middleware"
	"strings"

	"github.com/labstack/echo/v4"
)

// Helper: Mapping pesan default ke kode error
func parseErrorCode(message string, status int) string {
	msg := strings.ToLower(message)
	switch status {
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusUnauthorized:
		return "UNAUTHORIZED"
	case http.StatusForbidden:
		return "FORBIDDEN"
	case http.StatusBadRequest:
		if strings.Contains(msg, "validation") || strings.Contains(msg, "invalid") {
			return "VALIDATION_ERROR"
		}
		return "BAD_REQUEST"
	case http.StatusConflict:
		return "CONFLICT"
	default:
		return "INTERNAL_ERROR"
	}
}

// Custom global error handler
func CustomHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	msg := "Internal Server Error"
	var details interface{}

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code

		switch m := he.Message.(type) {
		case string:
			msg = m

		case dto.ErrorResponse:
			// ✅ handle your own custom error struct
			msg = m.Message
			details = m.Details

		case map[string]interface{}:
			if mMsg, ok := m["message"].(string); ok {
				msg = mMsg
			}
			details = m["details"]

		default:
			// fallback to avoid panic
			msg = fmt.Sprintf("%v", he.Message)
		}
	} else {
		msg = err.Error()
	}

	middleware.MakeLogEntry(c).Error(msg)
	errCode := parseErrorCode(msg, code)

	res := dto.ErrorResponse{
		Status:  code,
		Code:    errCode,
		Message: msg,
		Details: details,
	}

	if !c.Response().Committed {
		c.JSON(code, res)
	}
}
