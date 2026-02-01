package goyavision

import (
	"embed"
	"io/fs"
)

//go:embed all:web/dist
var webFS embed.FS

func GetWebFS() (fs.FS, error) {
	return fs.Sub(webFS, "web/dist")
}
