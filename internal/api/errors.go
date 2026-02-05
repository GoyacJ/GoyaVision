package api

import (
	"errors"
	"net/http"
	"strings"

	"goyavision/internal/api/dto"
	"goyavision/internal/api/response"
	"goyavision/pkg/apperr"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var (
	ErrNotFound           = errors.New("resource not found")
	ErrInvalidInput       = errors.New("invalid input")
	ErrAlreadyExists      = errors.New("resource already exists")
	ErrServiceUnavailable = errors.New("流媒体服务暂不可用")
)

func ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	if he, ok := err.(*echo.HTTPError); ok {
		_ = c.JSON(he.Code, dto.ErrorResponse{
			Error:   http.StatusText(he.Code),
			Message: he.Message.(string),
		})
		return
	}

	var appErr *apperr.Error
	if errors.As(err, &appErr) {
		status := mapErrorCodeToHTTPStatus(appErr.Code)
		_ = response.Err(c, status, appErr.Code, appErr.Message, appErr.Details)
		return
	}

	code := http.StatusInternalServerError
	message := err.Error()

	if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, ErrNotFound) {
		code = http.StatusNotFound
	} else if errors.Is(err, ErrInvalidInput) {
		code = http.StatusBadRequest
	} else if errors.Is(err, ErrAlreadyExists) {
		code = http.StatusConflict
	} else if strings.Contains(message, "mediamtx error") {
		code = http.StatusServiceUnavailable
		message = ErrServiceUnavailable.Error()
	}

	_ = c.JSON(code, dto.ErrorResponse{
		Error:   http.StatusText(code),
		Message: message,
	})
}

func mapErrorCodeToHTTPStatus(code int) int {
	switch {
	case code >= 40000 && code < 40100:
		return http.StatusBadRequest
	case code >= 40100 && code < 40200:
		return http.StatusUnauthorized
	case code >= 40300 && code < 40400:
		return http.StatusForbidden
	case code >= 40400 && code < 40500:
		return http.StatusNotFound
	case code >= 40900 && code < 41000:
		return http.StatusConflict
	case code >= 50300 && code < 50400:
		return http.StatusServiceUnavailable
	case code >= 50000 && code < 60000:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
