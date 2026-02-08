package dto

import (
	"time"

	"goyavision/internal/domain/algorithm"

	"github.com/google/uuid"
)

type AlgorithmListQuery struct {
	Status   *string `query:"status"`
	Scenario *string `query:"scenario"`
	Tags     *string `query:"tags"`
	Keyword  *string `query:"keyword"`
	Limit    int     `query:"limit"`
	Offset   int     `query:"offset"`
}

type AlgorithmImplementationReq struct {
	Name         string                 `json:"name,omitempty"`
	Type         string                 `json:"type,omitempty"`
	BindingRef   string                 `json:"binding_ref" validate:"required"`
	Config       map[string]interface{} `json:"config,omitempty"`
	LatencyMS    int                    `json:"latency_ms,omitempty"`
	CostScore    float64                `json:"cost_score,omitempty"`
	QualityScore float64                `json:"quality_score,omitempty"`
	Tier         string                 `json:"tier,omitempty"`
	IsDefault    bool                   `json:"is_default,omitempty"`
}

type AlgorithmEvaluationReq struct {
	DatasetRef       string             `json:"dataset_ref,omitempty"`
	Metrics          map[string]float64 `json:"metrics,omitempty"`
	ReportArtifactID *uuid.UUID         `json:"report_artifact_id,omitempty"`
	Summary          string             `json:"summary,omitempty"`
}

type AlgorithmVersionReq struct {
	Version         string                       `json:"version" validate:"required"`
	Status          string                       `json:"status,omitempty"`
	SelectionPolicy string                       `json:"selection_policy,omitempty"`
	Implementations []AlgorithmImplementationReq `json:"implementations"`
	Evaluations     []AlgorithmEvaluationReq     `json:"evaluations,omitempty"`
}

type AlgorithmCreateReq struct {
	Code           string               `json:"code" validate:"required"`
	Name           string               `json:"name" validate:"required"`
	Description    string               `json:"description,omitempty"`
	Scenario       string               `json:"scenario,omitempty"`
	Status         string               `json:"status,omitempty"`
	Tags           []string             `json:"tags,omitempty"`
	Visibility     *int                 `json:"visibility,omitempty"`
	VisibleRoleIDs []string             `json:"visible_role_ids,omitempty"`
	InitialVersion *AlgorithmVersionReq `json:"initial_version,omitempty"`
}

type AlgorithmUpdateReq struct {
	Name           *string  `json:"name,omitempty"`
	Description    *string  `json:"description,omitempty"`
	Scenario       *string  `json:"scenario,omitempty"`
	Status         *string  `json:"status,omitempty"`
	Tags           []string `json:"tags,omitempty"`
	Visibility     *int     `json:"visibility,omitempty"`
	VisibleRoleIDs []string `json:"visible_role_ids,omitempty"`
}

type CreateAlgorithmVersionReq struct {
	Version         string                       `json:"version" validate:"required"`
	Status          string                       `json:"status,omitempty"`
	SelectionPolicy string                       `json:"selection_policy,omitempty"`
	Implementations []AlgorithmImplementationReq `json:"implementations"`
	Evaluations     []AlgorithmEvaluationReq     `json:"evaluations,omitempty"`
}

type AlgorithmImplementationResponse struct {
	ID           uuid.UUID              `json:"id"`
	Name         string                 `json:"name"`
	Type         string                 `json:"type"`
	BindingRef   string                 `json:"binding_ref"`
	Config       map[string]interface{} `json:"config,omitempty"`
	LatencyMS    int                    `json:"latency_ms"`
	CostScore    float64                `json:"cost_score"`
	QualityScore float64                `json:"quality_score"`
	Tier         string                 `json:"tier"`
	IsDefault    bool                   `json:"is_default"`
}

type AlgorithmEvaluationResponse struct {
	ID               uuid.UUID          `json:"id"`
	DatasetRef       string             `json:"dataset_ref"`
	Metrics          map[string]float64 `json:"metrics,omitempty"`
	ReportArtifactID *uuid.UUID         `json:"report_artifact_id,omitempty"`
	Summary          string             `json:"summary,omitempty"`
}

type AlgorithmVersionResponse struct {
	ID                    uuid.UUID                         `json:"id"`
	Version               string                            `json:"version"`
	Status                string                            `json:"status"`
	SelectionPolicy       string                            `json:"selection_policy"`
	DefaultImplementation *uuid.UUID                        `json:"default_implementation,omitempty"`
	Implementations       []AlgorithmImplementationResponse `json:"implementations,omitempty"`
	Evaluations           []AlgorithmEvaluationResponse     `json:"evaluations,omitempty"`
}

type AlgorithmResponse struct {
	ID             uuid.UUID                  `json:"id"`
	Code           string                     `json:"code"`
	Name           string                     `json:"name"`
	Description    string                     `json:"description,omitempty"`
	Scenario       string                     `json:"scenario,omitempty"`
	Status         string                     `json:"status"`
	Tags           []string                   `json:"tags,omitempty"`
	Visibility     int                        `json:"visibility"`
	VisibleRoleIDs []string                   `json:"visible_role_ids,omitempty"`
	Versions       []AlgorithmVersionResponse `json:"versions,omitempty"`
	CreatedAt      time.Time                  `json:"created_at"`
	UpdatedAt      time.Time                  `json:"updated_at"`
}

type AlgorithmListResponse struct {
	Items []*AlgorithmResponse `json:"items"`
	Total int64                `json:"total"`
}

func AlgorithmToResponse(a *algorithm.Algorithm) *AlgorithmResponse {
	if a == nil {
		return nil
	}
	resp := &AlgorithmResponse{
		ID:             a.ID,
		Code:           a.Code,
		Name:           a.Name,
		Description:    a.Description,
		Scenario:       a.Scenario,
		Status:         string(a.Status),
		Tags:           a.Tags,
		Visibility:     int(a.Visibility),
		VisibleRoleIDs: a.VisibleRoleIDs,
		CreatedAt:      a.CreatedAt,
		UpdatedAt:      a.UpdatedAt,
	}
	if resp.Tags == nil {
		resp.Tags = []string{}
	}

	if len(a.Versions) > 0 {
		resp.Versions = make([]AlgorithmVersionResponse, 0, len(a.Versions))
		for i := range a.Versions {
			version := AlgorithmVersionToResponse(&a.Versions[i])
			if version != nil {
				resp.Versions = append(resp.Versions, *version)
			}
		}
	}

	return resp
}

func AlgorithmVersionToResponse(v *algorithm.Version) *AlgorithmVersionResponse {
	if v == nil {
		return nil
	}
	resp := &AlgorithmVersionResponse{
		ID:                    v.ID,
		Version:               v.Version,
		Status:                string(v.Status),
		SelectionPolicy:       string(v.SelectionPolicy),
		DefaultImplementation: v.DefaultImplementation,
	}
	for i := range v.Implementations {
		resp.Implementations = append(resp.Implementations, AlgorithmImplementationResponse{
			ID:           v.Implementations[i].ID,
			Name:         v.Implementations[i].Name,
			Type:         string(v.Implementations[i].Type),
			BindingRef:   v.Implementations[i].BindingRef,
			Config:       v.Implementations[i].Config,
			LatencyMS:    v.Implementations[i].LatencyMS,
			CostScore:    v.Implementations[i].CostScore,
			QualityScore: v.Implementations[i].QualityScore,
			Tier:         v.Implementations[i].Tier,
			IsDefault:    v.Implementations[i].IsDefault,
		})
	}
	for i := range v.Evaluations {
		resp.Evaluations = append(resp.Evaluations, AlgorithmEvaluationResponse{
			ID:               v.Evaluations[i].ID,
			DatasetRef:       v.Evaluations[i].DatasetRef,
			Metrics:          v.Evaluations[i].Metrics,
			ReportArtifactID: v.Evaluations[i].ReportArtifactID,
			Summary:          v.Evaluations[i].Summary,
		})
	}
	return resp
}

func AlgorithmsToResponse(items []*algorithm.Algorithm) []*AlgorithmResponse {
	out := make([]*AlgorithmResponse, len(items))
	for i := range items {
		out[i] = AlgorithmToResponse(items[i])
	}
	return out
}
