package handler

import (
	"goyavision/config"
	"goyavision/internal/adapter/mediamtx"
	"goyavision/internal/port"
)

type Deps struct {
	Repo   port.Repository
	Cfg    *config.Config
	MtxCli *mediamtx.Client
}
