package handler

import (
	"net/http"
	"strings"

	"goyavision/internal/api/dto"
	authmiddleware "goyavision/internal/api/middleware"
	appdto "goyavision/internal/app/dto"
	"goyavision/internal/domain/algorithm"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterAlgorithmRoutes(public *echo.Group, protected *echo.Group, h *Handlers) {
	handler := &algorithmHandler{h: h}

	public.GET("/algorithms", handler.List)
	public.GET("/algorithms/:id", handler.Get)

	protected.POST("/algorithms", handler.Create)
	protected.PUT("/algorithms/:id", handler.Update)
	protected.DELETE("/algorithms/:id", handler.Delete)
	protected.POST("/algorithms/:id/versions", handler.CreateVersion)
	protected.POST("/algorithms/:id/versions/:version_id/publish", handler.PublishVersion)
}

type algorithmHandler struct {
	h *Handlers
}

func (h *algorithmHandler) List(c echo.Context) error {
	var req dto.AlgorithmListQuery
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query parameters")
	}

	query := appdto.ListAlgorithmsQuery{
		Pagination: appdto.Pagination{
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	}
	if req.Status != nil && *req.Status != "" {
		status := algorithm.Status(*req.Status)
		query.Status = &status
	}
	if req.Scenario != nil {
		query.Scenario = *req.Scenario
	}
	if req.Tags != nil && *req.Tags != "" {
		query.Tags = strings.Split(*req.Tags, ",")
	}
	if req.Keyword != nil {
		query.Keyword = *req.Keyword
	}

	result, err := h.h.ListAlgorithms.Handle(c.Request().Context(), query)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.AlgorithmListResponse{
		Items: dto.AlgorithmsToResponse(result.Items),
		Total: result.Total,
	})
}

func (h *algorithmHandler) Get(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid algorithm id")
	}
	result, err := h.h.GetAlgorithm.Handle(c.Request().Context(), appdto.GetAlgorithmQuery{ID: id})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, dto.AlgorithmToResponse(result))
}

func (h *algorithmHandler) Create(c echo.Context) error {
	var req dto.AlgorithmCreateReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	status := algorithm.Status(req.Status)
	if status == "" {
		status = algorithm.StatusDraft
	}
	visibility := algorithm.VisibilityPrivate
	if req.Visibility != nil {
		visibility = algorithm.Visibility(*req.Visibility)
	}

	visibleRoleIDs := req.VisibleRoleIDs
	if visibility == algorithm.VisibilityRole && len(visibleRoleIDs) == 0 {
		if ids, ok := c.Get(authmiddleware.ContextKeyRoleIDs).([]uuid.UUID); ok {
			for _, id := range ids {
				visibleRoleIDs = append(visibleRoleIDs, id.String())
			}
		}
	}

	cmd := appdto.CreateAlgorithmCommand{
		Code:           req.Code,
		Name:           req.Name,
		Description:    req.Description,
		Scenario:       req.Scenario,
		Status:         status,
		Tags:           req.Tags,
		Visibility:     visibility,
		VisibleRoleIDs: visibleRoleIDs,
		InitialVersion: convertAlgorithmVersionInput(req.InitialVersion),
	}

	result, err := h.h.CreateAlgorithm.Handle(c.Request().Context(), cmd)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, dto.AlgorithmToResponse(result))
}

func (h *algorithmHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid algorithm id")
	}

	var req dto.AlgorithmUpdateReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	cmd := appdto.UpdateAlgorithmCommand{
		ID:             id,
		Name:           req.Name,
		Description:    req.Description,
		Scenario:       req.Scenario,
		Tags:           req.Tags,
		VisibleRoleIDs: req.VisibleRoleIDs,
	}
	if req.Status != nil {
		status := algorithm.Status(*req.Status)
		cmd.Status = &status
	}
	if req.Visibility != nil {
		visibility := algorithm.Visibility(*req.Visibility)
		cmd.Visibility = &visibility
	}

	result, err := h.h.UpdateAlgorithm.Handle(c.Request().Context(), cmd)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, dto.AlgorithmToResponse(result))
}

func (h *algorithmHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid algorithm id")
	}
	if err := h.h.DeleteAlgorithm.Handle(c.Request().Context(), appdto.DeleteAlgorithmCommand{ID: id}); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *algorithmHandler) CreateVersion(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid algorithm id")
	}

	var req dto.CreateAlgorithmVersionReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	cmd := appdto.CreateAlgorithmVersionCommand{
		AlgorithmID:     id,
		Version:         req.Version,
		Status:          algorithm.VersionStatus(req.Status),
		SelectionPolicy: algorithm.SelectionPolicy(req.SelectionPolicy),
		Implementations: convertAlgorithmImplementationInputs(req.Implementations),
		Evaluations:     convertAlgorithmEvaluationInputs(req.Evaluations),
	}

	version, err := h.h.CreateAlgorithmVersion.Handle(c.Request().Context(), cmd)
	if err != nil {
		return err
	}
	resp := dto.AlgorithmVersionToResponse(version)
	if resp == nil {
		return c.JSON(http.StatusCreated, map[string]interface{}{})
	}
	return c.JSON(http.StatusCreated, resp)
}

func (h *algorithmHandler) PublishVersion(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid algorithm id")
	}
	versionID, err := uuid.Parse(c.Param("version_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid version id")
	}

	version, err := h.h.PublishAlgorithmVersion.Handle(c.Request().Context(), appdto.PublishAlgorithmVersionCommand{
		AlgorithmID: id,
		VersionID:   versionID,
	})
	if err != nil {
		return err
	}
	resp := dto.AlgorithmVersionToResponse(version)
	if resp == nil {
		return c.JSON(http.StatusOK, map[string]interface{}{})
	}
	return c.JSON(http.StatusOK, resp)
}

func convertAlgorithmVersionInput(req *dto.AlgorithmVersionReq) *appdto.AlgorithmVersionInput {
	if req == nil {
		return nil
	}
	return &appdto.AlgorithmVersionInput{
		Version:         req.Version,
		Status:          algorithm.VersionStatus(req.Status),
		SelectionPolicy: algorithm.SelectionPolicy(req.SelectionPolicy),
		Implementations: convertAlgorithmImplementationInputs(req.Implementations),
		Evaluations:     convertAlgorithmEvaluationInputs(req.Evaluations),
	}
}

func convertAlgorithmImplementationInputs(reqs []dto.AlgorithmImplementationReq) []appdto.AlgorithmImplementationInput {
	if len(reqs) == 0 {
		return nil
	}
	out := make([]appdto.AlgorithmImplementationInput, 0, len(reqs))
	for i := range reqs {
		out = append(out, appdto.AlgorithmImplementationInput{
			Name:         reqs[i].Name,
			Type:         algorithm.ImplementationType(reqs[i].Type),
			BindingRef:   reqs[i].BindingRef,
			Config:       reqs[i].Config,
			LatencyMS:    reqs[i].LatencyMS,
			CostScore:    reqs[i].CostScore,
			QualityScore: reqs[i].QualityScore,
			Tier:         reqs[i].Tier,
			IsDefault:    reqs[i].IsDefault,
		})
	}
	return out
}

func convertAlgorithmEvaluationInputs(reqs []dto.AlgorithmEvaluationReq) []appdto.AlgorithmEvaluationInput {
	if len(reqs) == 0 {
		return nil
	}
	out := make([]appdto.AlgorithmEvaluationInput, 0, len(reqs))
	for i := range reqs {
		out = append(out, appdto.AlgorithmEvaluationInput{
			DatasetRef:       reqs[i].DatasetRef,
			Metrics:          reqs[i].Metrics,
			ReportArtifactID: reqs[i].ReportArtifactID,
			Summary:          reqs[i].Summary,
		})
	}
	return out
}
