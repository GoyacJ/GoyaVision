package api

import (
	"log"

	"github.com/labstack/echo/v4"
)

func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			c.Error(err)
			return err
		}
		log.Printf("%s %s %d", c.Request().Method, c.Path(), c.Response().Status)
		return nil
	}
}

func Recover(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic: %v", r)
				err = c.JSON(500, map[string]string{"error": "internal server error"})
			}
		}()
		return next(c)
	}
}
