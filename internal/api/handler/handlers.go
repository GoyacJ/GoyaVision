package handler

import (
	"goyavision/config"
	"goyavision/internal/adapter/mediamtx"
	"goyavision/internal/app"
	"goyavision/internal/app/command"
	appport "goyavision/internal/app/port"
	"goyavision/internal/app/query"
	"goyavision/internal/port"
	"goyavision/pkg/storage"
)

type Handlers struct {
	CreateSource         *command.CreateSourceHandler
	UpdateSource         *command.UpdateSourceHandler
	DeleteSource         *command.DeleteSourceHandler
	CreateAsset          *command.CreateAssetHandler
	UpdateAsset          *command.UpdateAssetHandler
	DeleteAsset          *command.DeleteAssetHandler
	Login                *command.LoginHandler
	CreateOperator       *command.CreateOperatorHandler
	UpdateOperator       *command.UpdateOperatorHandler
	DeleteOperator       *command.DeleteOperatorHandler
	EnableOperator       *command.EnableOperatorHandler
	CreateWorkflow       *command.CreateWorkflowHandler
	UpdateWorkflow       *command.UpdateWorkflowHandler
	DeleteWorkflow       *command.DeleteWorkflowHandler
	EnableWorkflow       *command.EnableWorkflowHandler
	CreateTask           *command.CreateTaskHandler
	UpdateTask           *command.UpdateTaskHandler
	DeleteTask           *command.DeleteTaskHandler
	StartTask            *command.StartTaskHandler
	CompleteTask         *command.CompleteTaskHandler
	FailTask             *command.FailTaskHandler
	CancelTask           *command.CancelTaskHandler
	GetSource            *query.GetSourceHandler
	ListSources          *query.ListSourcesHandler
	GetAsset             *query.GetAssetHandler
	ListAssets           *query.ListAssetsHandler
	ListAssetChildren    *query.ListAssetChildrenHandler
	GetAssetTags         *query.GetAssetTagsHandler
	GetProfile           *query.GetProfileHandler
	GetOperator          *query.GetOperatorHandler
	GetOperatorByCode    *query.GetOperatorByCodeHandler
	ListOperators        *query.ListOperatorsHandler
	GetWorkflow          *query.GetWorkflowHandler
	GetWorkflowWithNodes *query.GetWorkflowWithNodesHandler
	GetWorkflowByCode    *query.GetWorkflowByCodeHandler
	ListWorkflows        *query.ListWorkflowsHandler
	GetTask              *query.GetTaskHandler
	GetTaskWithRelations *query.GetTaskWithRelationsHandler
	ListTasks            *query.ListTasksHandler
	GetTaskStats         *query.GetTaskStatsHandler
	ListRunningTasks     *query.ListRunningTasksHandler
	Cfg                  *config.Config
	MtxCli               *mediamtx.Client
	MinIOClient          *storage.MinIOClient
	WorkflowScheduler    *app.WorkflowScheduler
	Repo                 port.Repository // For middleware and non-migrated handlers
	TokenService         appport.TokenService // For auth handlers
}

// Deps 依赖注入结构
type Deps struct {
	Repo        port.Repository
	Cfg         *config.Config
	MinIOClient *storage.MinIOClient
}

func NewHandlers(
	uow appport.UnitOfWork,
	mediaGateway appport.MediaGateway,
	tokenService appport.TokenService,
	cfg *config.Config,
	mtxCli *mediamtx.Client,
	minioClient *storage.MinIOClient,
	workflowScheduler *app.WorkflowScheduler,
	repo port.Repository,
) *Handlers {
	return &Handlers{
		CreateSource:         command.NewCreateSourceHandler(uow, mediaGateway),
		UpdateSource:         command.NewUpdateSourceHandler(uow, mediaGateway),
		DeleteSource:         command.NewDeleteSourceHandler(uow, mediaGateway),
		CreateAsset:          command.NewCreateAssetHandler(uow),
		UpdateAsset:          command.NewUpdateAssetHandler(uow),
		DeleteAsset:          command.NewDeleteAssetHandler(uow),
		Login:                command.NewLoginHandler(uow, tokenService),
		CreateOperator:       command.NewCreateOperatorHandler(uow),
		UpdateOperator:       command.NewUpdateOperatorHandler(uow),
		DeleteOperator:       command.NewDeleteOperatorHandler(uow),
		EnableOperator:       command.NewEnableOperatorHandler(uow),
		CreateWorkflow:       command.NewCreateWorkflowHandler(uow),
		UpdateWorkflow:       command.NewUpdateWorkflowHandler(uow),
		DeleteWorkflow:       command.NewDeleteWorkflowHandler(uow),
		EnableWorkflow:       command.NewEnableWorkflowHandler(uow),
		CreateTask:           command.NewCreateTaskHandler(uow),
		UpdateTask:           command.NewUpdateTaskHandler(uow),
		DeleteTask:           command.NewDeleteTaskHandler(uow),
		StartTask:            command.NewStartTaskHandler(uow),
		CompleteTask:         command.NewCompleteTaskHandler(uow),
		FailTask:             command.NewFailTaskHandler(uow),
		CancelTask:           command.NewCancelTaskHandler(uow),
		GetSource:            query.NewGetSourceHandler(uow),
		ListSources:          query.NewListSourcesHandler(uow),
		GetAsset:             query.NewGetAssetHandler(uow),
		ListAssets:           query.NewListAssetsHandler(uow),
		ListAssetChildren:    query.NewListAssetChildrenHandler(uow),
		GetAssetTags:         query.NewGetAssetTagsHandler(uow),
		GetProfile:           query.NewGetProfileHandler(uow),
		GetOperator:          query.NewGetOperatorHandler(uow),
		GetOperatorByCode:    query.NewGetOperatorByCodeHandler(uow),
		ListOperators:        query.NewListOperatorsHandler(uow),
		GetWorkflow:          query.NewGetWorkflowHandler(uow),
		GetWorkflowWithNodes: query.NewGetWorkflowWithNodesHandler(uow),
		GetWorkflowByCode:    query.NewGetWorkflowByCodeHandler(uow),
		ListWorkflows:        query.NewListWorkflowsHandler(uow),
		GetTask:              query.NewGetTaskHandler(uow),
		GetTaskWithRelations: query.NewGetTaskWithRelationsHandler(uow),
		ListTasks:            query.NewListTasksHandler(uow),
		GetTaskStats:         query.NewGetTaskStatsHandler(uow),
		ListRunningTasks:     query.NewListRunningTasksHandler(uow),
		Cfg:                  cfg,
		MtxCli:               mtxCli,
		MinIOClient:          minioClient,
		WorkflowScheduler:    workflowScheduler,
		Repo:                 repo,
		TokenService:         tokenService,
	}
}
