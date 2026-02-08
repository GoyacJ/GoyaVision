package handler

import (
	"goyavision/config"
	"goyavision/internal/adapter/crypto"
	"goyavision/internal/adapter/engine"
	"goyavision/internal/adapter/payment"
	"goyavision/internal/adapter/mediamtx"
	"goyavision/internal/app"
	"goyavision/internal/app/command"
	appport "goyavision/internal/app/port"
	"goyavision/internal/app/query"
	infraauth "goyavision/internal/infra/auth"
	"goyavision/internal/port"
	"gorm.io/gorm"
)

type Handlers struct {
	CreateSource             *command.CreateSourceHandler
	UpdateSource             *command.UpdateSourceHandler
	DeleteSource             *command.DeleteSourceHandler
	CreateAsset              *command.CreateAssetHandler
	UpdateAsset              *command.UpdateAssetHandler
	DeleteAsset              *command.DeleteAssetHandler
	Login                    *command.LoginHandler
	LoginOAuth               *command.LoginOAuthHandler
	BindIdentity             *command.BindIdentityHandler
	CreateOperator           *command.CreateOperatorHandler
	UpdateOperator           *command.UpdateOperatorHandler
	DeleteOperator           *command.DeleteOperatorHandler
	CreateOperatorVersion    *command.CreateOperatorVersionHandler
	ActivateVersion          *command.ActivateVersionHandler
	RollbackVersion          *command.RollbackVersionHandler
	ArchiveVersion           *command.ArchiveVersionHandler
	InstallTemplate          *command.InstallTemplateHandler
	SetOperatorDependencies  *command.SetOperatorDependenciesHandler
	PublishOperator          *command.PublishOperatorHandler
	DeprecateOperator        *command.DeprecateOperatorHandler
	TestOperator             *command.TestOperatorHandler
	InstallMCPOperator       *command.InstallMCPOperatorHandler
	SyncMCPTemplates         *command.SyncMCPTemplatesHandler
	CreateAIModel            *command.CreateAIModelHandler
	UpdateAIModel            *command.UpdateAIModelHandler
	DeleteAIModel            *command.DeleteAIModelHandler
	TestAIModel              *command.TestAIModelHandler
	ListAIModels             *query.ListAIModelsHandler
	GetAIModel               *query.GetAIModelHandler
	CreateWorkflow           *command.CreateWorkflowHandler
	UpdateWorkflow           *command.UpdateWorkflowHandler
	DeleteWorkflow           *command.DeleteWorkflowHandler
	EnableWorkflow           *command.EnableWorkflowHandler
	CreateTask               *command.CreateTaskHandler
	UpdateTask               *command.UpdateTaskHandler
	DeleteTask               *command.DeleteTaskHandler
	StartTask                *command.StartTaskHandler
	CompleteTask             *command.CompleteTaskHandler
	FailTask                 *command.FailTaskHandler
	CancelTask               *command.CancelTaskHandler
	GetSource                *query.GetSourceHandler
	ListSources              *query.ListSourcesHandler
	GetAsset                 *query.GetAssetHandler
	ListAssets               *query.ListAssetsHandler
	ListAssetChildren        *query.ListAssetChildrenHandler
	GetAssetTags             *query.GetAssetTagsHandler
	GetProfile               *query.GetProfileHandler
	GetOperator              *query.GetOperatorHandler
	GetOperatorByCode        *query.GetOperatorByCodeHandler
	ListOperators            *query.ListOperatorsHandler
	ListOperatorVersions     *query.ListOperatorVersionsHandler
	GetOperatorVersion       *query.GetOperatorVersionHandler
	ListTemplates            *query.ListTemplatesHandler
	GetTemplate              *query.GetTemplateHandler
	ListOperatorDependencies *query.ListOperatorDependenciesHandler
	CheckDependencies        *query.CheckDependenciesHandler
	ValidateSchema           *query.ValidateSchemaHandler
	ValidateConnection       *query.ValidateConnectionHandler
	ListMCPServers           *query.ListMCPServersHandler
	ListMCPTools             *query.ListMCPToolsHandler
	PreviewMCPTool           *query.PreviewMCPToolHandler
	GetWorkflow              *query.GetWorkflowHandler
	GetWorkflowWithNodes     *query.GetWorkflowWithNodesHandler
	GetWorkflowByCode        *query.GetWorkflowByCodeHandler
	ListWorkflows            *query.ListWorkflowsHandler
	GetTask                  *query.GetTaskHandler
	GetTaskWithRelations     *query.GetTaskWithRelationsHandler
	ListTasks                *query.ListTasksHandler
	GetTaskStats             *query.GetTaskStatsHandler
	ListRunningTasks         *query.ListRunningTasksHandler
	GetUserAssetSummaryHandler      *query.GetUserAssetSummaryHandler
	ListUserTransactionsHandler     *query.ListUserTransactionsHandler
	ListUserPointRecordsHandler     *query.ListUserPointRecordsHandler
	GetUserUsageStatsHandler        *query.GetUserUsageStatsHandler
	RechargeHandler                 *command.RechargeHandler
	CheckInHandler                  *command.CheckInHandler
	SubscribeHandler                *command.SubscribeHandler
	Cfg               *config.Config
	MtxCli            *mediamtx.Client
	FileStorage       appport.FileStorage
	StorageURLConfig  appport.StorageURLConfig
	WorkflowScheduler *app.WorkflowScheduler
	DB                       *gorm.DB
	Repo                     port.Repository      // For middleware and non-migrated handlers
	TokenService             appport.TokenService // For auth handlers
	AuthProviderFactory      appport.AuthProviderFactory
}

// Deps 依赖注入结构
type Deps struct {
	Repo             port.Repository
	Cfg              *config.Config
	FileStorage      appport.FileStorage
	StorageURLConfig appport.StorageURLConfig
}

func NewHandlers(
	uow appport.UnitOfWork,
	schemaValidator appport.SchemaValidator,
	mcpClient port.MCPClient,
	mcpRegistry port.MCPRegistry,
	mediaGateway appport.MediaGateway,
	tokenService appport.TokenService,
	db *gorm.DB,
	cfg *config.Config,
	mtxCli *mediamtx.Client,
	fileStorage appport.FileStorage,
	storageURLConfig appport.StorageURLConfig,
	workflowScheduler *app.WorkflowScheduler,
	repo port.Repository,
	eventBus appport.EventBus,
) *Handlers {
	authProviderFactory := infraauth.NewProviderFactory(cfg)
	userService := app.NewUserService(repo)
	encryptKey := cfg.EncryptKey
	if encryptKey == "" {
		encryptKey = cfg.JWT.Secret
	}
	cryptoService, _ := crypto.NewAESCryptoService(encryptKey)

	httpExecutor := engine.NewHTTPOperatorExecutor()
	cliExecutor := engine.NewCLIOperatorExecutor()
	mcpExecutor := engine.NewMCPOperatorExecutor(mcpClient)
	aiModelExecutor := engine.NewAIModelExecutor(repo, cryptoService)
	executorRegistry := engine.NewExecutorRegistry()
	executorRegistry.Register(httpExecutor.Mode(), httpExecutor)
	executorRegistry.Register(cliExecutor.Mode(), cliExecutor)
	executorRegistry.Register(mcpExecutor.Mode(), mcpExecutor)
	executorRegistry.Register(aiModelExecutor.Mode(), aiModelExecutor)

	paymentAdapter, _ := payment.NewGoPayAdapter(cfg.Payment)

	return &Handlers{
		CreateSource:             command.NewCreateSourceHandler(uow, mediaGateway),
		UpdateSource:             command.NewUpdateSourceHandler(uow, mediaGateway),
		DeleteSource:             command.NewDeleteSourceHandler(uow, mediaGateway),
		CreateAsset:              command.NewCreateAssetHandler(uow, eventBus),
		UpdateAsset:              command.NewUpdateAssetHandler(uow),
		DeleteAsset:              command.NewDeleteAssetHandler(uow),
		Login:                    command.NewLoginHandler(uow, tokenService),
		LoginOAuth:               command.NewLoginOAuthHandler(uow, tokenService, authProviderFactory, userService),
		BindIdentity:             command.NewBindIdentityHandler(uow, tokenService),
		CreateOperator:           command.NewCreateOperatorHandler(uow, schemaValidator),
		UpdateOperator:           command.NewUpdateOperatorHandler(uow),
		DeleteOperator:           command.NewDeleteOperatorHandler(uow),
		CreateOperatorVersion:    command.NewCreateOperatorVersionHandler(uow, schemaValidator),
		ActivateVersion:          command.NewActivateVersionHandler(uow),
		RollbackVersion:          command.NewRollbackVersionHandler(uow),
		ArchiveVersion:           command.NewArchiveVersionHandler(uow),
		InstallTemplate:          command.NewInstallTemplateHandler(uow),
		SetOperatorDependencies:  command.NewSetOperatorDependenciesHandler(uow),
		PublishOperator:          command.NewPublishOperatorHandler(uow, mcpClient, schemaValidator),
		DeprecateOperator:        command.NewDeprecateOperatorHandler(uow),
		TestOperator:             command.NewTestOperatorHandler(uow, executorRegistry),
		InstallMCPOperator:       command.NewInstallMCPOperatorHandler(uow, mcpClient),
		SyncMCPTemplates:         command.NewSyncMCPTemplatesHandler(uow, mcpClient),
		CreateAIModel:            command.NewCreateAIModelHandler(uow, cryptoService),
		UpdateAIModel:            command.NewUpdateAIModelHandler(uow, cryptoService),
		DeleteAIModel:            command.NewDeleteAIModelHandler(uow),
		TestAIModel:              command.NewTestAIModelHandler(repo, cryptoService),
		ListAIModels:             query.NewListAIModelsHandler(repo),
		GetAIModel:               query.NewGetAIModelHandler(repo),
		CreateWorkflow:           command.NewCreateWorkflowHandler(uow, schemaValidator),
		UpdateWorkflow:           command.NewUpdateWorkflowHandler(uow, schemaValidator),
		DeleteWorkflow:           command.NewDeleteWorkflowHandler(uow),
		EnableWorkflow:           command.NewEnableWorkflowHandler(uow),
		CreateTask:               command.NewCreateTaskHandler(uow),
		UpdateTask:               command.NewUpdateTaskHandler(uow),
		DeleteTask:               command.NewDeleteTaskHandler(uow),
		StartTask:                command.NewStartTaskHandler(uow),
		CompleteTask:             command.NewCompleteTaskHandler(uow),
		FailTask:                 command.NewFailTaskHandler(uow),
		CancelTask:               command.NewCancelTaskHandler(uow),
		GetSource:                query.NewGetSourceHandler(uow),
		ListSources:              query.NewListSourcesHandler(uow),
		GetAsset:                 query.NewGetAssetHandler(uow),
		ListAssets:               query.NewListAssetsHandler(uow),
		ListAssetChildren:        query.NewListAssetChildrenHandler(uow),
		GetAssetTags:             query.NewGetAssetTagsHandler(uow),
		GetProfile:               query.NewGetProfileHandler(uow),
		GetOperator:              query.NewGetOperatorHandler(uow),
		GetOperatorByCode:        query.NewGetOperatorByCodeHandler(uow),
		ListOperators:            query.NewListOperatorsHandler(uow),
		ListOperatorVersions:     query.NewListOperatorVersionsHandler(uow),
		GetOperatorVersion:       query.NewGetOperatorVersionHandler(uow),
		ListTemplates:            query.NewListTemplatesHandler(uow),
		GetTemplate:              query.NewGetTemplateHandler(uow),
		ListOperatorDependencies: query.NewListOperatorDependenciesHandler(uow),
		CheckDependencies:        query.NewCheckDependenciesHandler(uow),
		ValidateSchema:           query.NewValidateSchemaHandler(schemaValidator),
		ValidateConnection:       query.NewValidateConnectionHandler(schemaValidator),
		ListMCPServers:           query.NewListMCPServersHandler(mcpRegistry),
		ListMCPTools:             query.NewListMCPToolsHandler(mcpClient),
		PreviewMCPTool:           query.NewPreviewMCPToolHandler(mcpClient),
		GetWorkflow:              query.NewGetWorkflowHandler(uow),
		GetWorkflowWithNodes:     query.NewGetWorkflowWithNodesHandler(uow),
		GetWorkflowByCode:        query.NewGetWorkflowByCodeHandler(uow),
		ListWorkflows:            query.NewListWorkflowsHandler(uow),
		GetTask:                  query.NewGetTaskHandler(uow),
		GetTaskWithRelations:     query.NewGetTaskWithRelationsHandler(uow),
		ListTasks:                query.NewListTasksHandler(uow),
		GetTaskStats:             query.NewGetTaskStatsHandler(uow),
		ListRunningTasks:         query.NewListRunningTasksHandler(uow),
		GetUserAssetSummaryHandler:      query.NewGetUserAssetSummaryHandler(uow),
		ListUserTransactionsHandler:     query.NewListUserTransactionsHandler(uow),
		ListUserPointRecordsHandler:     query.NewListUserPointRecordsHandler(uow),
		GetUserUsageStatsHandler:        query.NewGetUserUsageStatsHandler(uow),
		RechargeHandler:                 command.NewRechargeHandler(uow, paymentAdapter),
		CheckInHandler:                  command.NewCheckInHandler(uow),
		SubscribeHandler:                command.NewSubscribeHandler(uow),
		Cfg:               cfg,
		MtxCli:            mtxCli,
		FileStorage:       fileStorage,
		StorageURLConfig:  storageURLConfig,
		WorkflowScheduler: workflowScheduler,
		DB:                       db,
		Repo:                     repo,
		TokenService:             tokenService,
		AuthProviderFactory:      authProviderFactory,
	}
}
