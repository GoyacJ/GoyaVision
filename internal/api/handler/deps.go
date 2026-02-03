package handler

import (
	"goyavision/config"
	"goyavision/internal/adapter/mediamtx"
	"goyavision/internal/app"
	"goyavision/internal/port"
	"goyavision/pkg/storage"
)

type Deps struct {
	Repo                port.Repository
	Cfg                 *config.Config
	MtxCli              *mediamtx.Client
	MediaSourceService  *app.MediaSourceService
	MinIOClient         *storage.MinIOClient
	WorkflowScheduler   *app.WorkflowScheduler
}
