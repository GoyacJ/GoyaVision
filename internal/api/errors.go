package api

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var (
	ErrNotFound      = errors.New("resource not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrAlreadyExists = errors.New("resource already exists")
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func ErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := err.Error()

	if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, ErrNotFound) {
		code = http.StatusNotFound
	} else if errors.Is(err, ErrInvalidInput) {
		code = http.StatusBadRequest
	} else if errors.Is(err, ErrAlreadyExists) {
		code = http.StatusConflict
	}

	if !c.Response().Committed {
		c.JSON(code, ErrorResponse{
			Error:   http.StatusText(code),
			Message: message,
		})
	}
}
