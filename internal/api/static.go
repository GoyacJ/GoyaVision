package api

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

var (
	//go:embed all:web/dist
	webFS embed.FS
)

func RegisterStatic(e *echo.Echo) {
	webDist, err := fs.Sub(webFS, "web/dist")
	if err != nil {
		return
	}

	fileServer := http.FileServer(http.FS(webDist))

	e.GET("/*", func(c echo.Context) error {
		path := c.Request().URL.Path

		if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/live/") {
			return c.NoContent(http.StatusNotFound)
		}

		if path == "/" || path == "" {
			path = "/index.html"
		}

		file, err := webDist.Open(strings.TrimPrefix(path, "/"))
		if err != nil {
			indexFile, err := webDist.Open("index.html")
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
			indexFile, err := webDist.Open("index.html")
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
	switch {
	case path[len(path)-5:] == ".html":
		return "text/html"
	case path[len(path)-4:] == ".css":
		return "text/css"
	case path[len(path)-3:] == ".js":
		return "application/javascript"
	case path[len(path)-4:] == ".json":
		return "application/json"
	case path[len(path)-4:] == ".png":
		return "image/png"
	case path[len(path)-4:] == ".jpg" || path[len(path)-5:] == ".jpeg":
		return "image/jpeg"
	case path[len(path)-4:] == ".svg":
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}
