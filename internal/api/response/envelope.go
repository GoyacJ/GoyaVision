package response

import (
	"net/http"
	"time"

	"goyavision/pkg/logger"

	"github.com/labstack/echo/v4"
)

type Envelope struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Details   interface{} `json:"details,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

type PagedData struct {
	Items    interface{} `json:"items"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

func newEnvelope(c echo.Context, code int, message string, data interface{}) Envelope {
	return Envelope{
		Code:      code,
		Message:   message,
		Data:      data,
		RequestID: logger.RequestIDFromContext(c.Request().Context()),
		Timestamp: time.Now().Unix(),
	}
}

func OK(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, newEnvelope(c, 0, "success", data))
}

func Created(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, newEnvelope(c, 0, "created", data))
}

func Paged(c echo.Context, items interface{}, total int64, page, pageSize int) error {
	data := PagedData{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
	return c.JSON(http.StatusOK, newEnvelope(c, 0, "success", data))
}

func NoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func Err(c echo.Context, httpStatus int, code int, message string, details interface{}) error {
	env := Envelope{
		Code:      code,
		Message:   message,
		Details:   details,
		RequestID: logger.RequestIDFromContext(c.Request().Context()),
		Timestamp: time.Now().Unix(),
	}
	return c.JSON(httpStatus, env)
}
