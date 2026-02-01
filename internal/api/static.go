package api

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func RegisterStatic(e *echo.Echo, webFS fs.FS) {
	if webFS == nil {
		return
	}

	e.GET("/*", func(c echo.Context) error {
		path := c.Request().URL.Path

		if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/live/") {
			return c.NoContent(http.StatusNotFound)
		}

		if path == "/" || path == "" {
			path = "/index.html"
		}

		file, err := webFS.Open(strings.TrimPrefix(path, "/"))
		if err != nil {
			indexFile, err := webFS.Open("index.html")
			if err != nil {
				return c.String(http.StatusNotFound, "Not Found")
			}
			defer indexFile.Close()
			return c.Stream(http.StatusOK, "text/html", indexFile)
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			return err
		}

		if stat.IsDir() {
			indexFile, err := webFS.Open("index.html")
			if err != nil {
				return c.String(http.StatusNotFound, "Not Found")
			}
			defer indexFile.Close()
			return c.Stream(http.StatusOK, "text/html", indexFile)
		}

		return c.Stream(http.StatusOK, getContentType(path), file)
	})
}

func getContentType(path string) string {
	if len(path) < 3 {
		return "application/octet-stream"
	}
	switch {
	case strings.HasSuffix(path, ".html"):
		return "text/html"
	case strings.HasSuffix(path, ".css"):
		return "text/css"
	case strings.HasSuffix(path, ".js"):
		return "application/javascript"
	case strings.HasSuffix(path, ".json"):
		return "application/json"
	case strings.HasSuffix(path, ".png"):
		return "image/png"
	case strings.HasSuffix(path, ".jpg"), strings.HasSuffix(path, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(path, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(path, ".woff2"):
		return "font/woff2"
	case strings.HasSuffix(path, ".woff"):
		return "font/woff"
	default:
		return "application/octet-stream"
	}
}
