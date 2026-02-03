package api

import (
	"errors"
	"net/http"
	"strings"

	"goyavision/internal/api/dto"
	"goyavision/internal/app"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var (
	ErrNotFound         = errors.New("resource not found")
	ErrInvalidInput     = errors.New("invalid input")
	ErrAlreadyExists    = errors.New("resource already exists")
	ErrServiceUnavailable = errors.New("流媒体服务暂不可用")
)

func ErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := err.Error()

	if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, ErrNotFound) {
		code = http.StatusNotFound
	} else if errors.Is(err, ErrInvalidInput) {
		code = http.StatusBadRequest
	} else if errors.Is(err, ErrAlreadyExists) || errors.Is(err, app.ErrMediaSourceHasAssets) {
		code = http.StatusConflict
	} else if errors.Is(err, app.ErrMediaMTXUnavailable) || strings.Contains(message, "mediamtx error") {
		code = http.StatusServiceUnavailable
		if errors.Is(err, app.ErrMediaMTXUnavailable) {
			message = ErrServiceUnavailable.Error()
		}
	}

	if !c.Response().Committed {
		c.JSON(code, dto.ErrorResponse{
			Error:   http.StatusText(code),
			Message: message,
		})
	}
}
