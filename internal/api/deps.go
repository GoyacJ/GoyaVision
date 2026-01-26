package api

import (
	"goyavision/config"
	"goyavision/internal/port"
)

type Deps struct {
	Repo port.Repository
	Cfg  *config.Config
}
