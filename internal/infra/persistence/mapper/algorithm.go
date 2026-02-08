package mapper

import (
	"encoding/json"

	"goyavision/internal/domain/algorithm"
	"goyavision/internal/infra/persistence/model"

	"gorm.io/datatypes"
)

func AlgorithmToModel(a *algorithm.Algorithm) *model.AlgorithmModel {
	m := &model.AlgorithmModel{
		ID:          a.ID,
		TenantID:    a.TenantID,
		OwnerID:     a.OwnerID,
		Visibility:  int(a.Visibility),
		Code:        a.Code,
		Name:        a.Name,
		Description: a.Description,
		Scenario:    a.Scenario,
		Status:      string(a.Status),
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
	if a.VisibleRoleIDs != nil {
		b, _ := json.Marshal(a.VisibleRoleIDs)
		m.VisibleRoleIDs = datatypes.JSON(b)
	}
	if a.Tags != nil {
		b, _ := json.Marshal(a.Tags)
		m.Tags = datatypes.JSON(b)
	}
	for i := range a.Versions {
		m.Versions = append(m.Versions, *AlgorithmVersionToModel(&a.Versions[i]))
	}
	return m
}

func AlgorithmToDomain(m *model.AlgorithmModel) *algorithm.Algorithm {
	a := &algorithm.Algorithm{
		ID:          m.ID,
		TenantID:    m.TenantID,
		OwnerID:     m.OwnerID,
		Visibility:  algorithm.Visibility(m.Visibility),
		Code:        m.Code,
		Name:        m.Name,
		Description: m.Description,
		Scenario:    m.Scenario,
		Status:      algorithm.Status(m.Status),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
	if m.VisibleRoleIDs != nil {
		_ = json.Unmarshal(m.VisibleRoleIDs, &a.VisibleRoleIDs)
	}
	if m.Tags != nil {
		_ = json.Unmarshal(m.Tags, &a.Tags)
	}
	for i := range m.Versions {
		a.Versions = append(a.Versions, *AlgorithmVersionToDomain(&m.Versions[i]))
	}
	return a
}

func AlgorithmVersionToModel(v *algorithm.Version) *model.AlgorithmVersionModel {
	m := &model.AlgorithmVersionModel{
		ID:                      v.ID,
		AlgorithmID:             v.AlgorithmID,
		Version:                 v.Version,
		Status:                  string(v.Status),
		SelectionPolicy:         string(v.SelectionPolicy),
		DefaultImplementationID: v.DefaultImplementation,
		CreatedAt:               v.CreatedAt,
		UpdatedAt:               v.UpdatedAt,
	}
	for i := range v.Implementations {
		m.Implementations = append(m.Implementations, *AlgorithmImplementationToModel(&v.Implementations[i]))
	}
	for i := range v.Evaluations {
		m.Evaluations = append(m.Evaluations, *AlgorithmEvaluationToModel(&v.Evaluations[i]))
	}
	return m
}

func AlgorithmVersionToDomain(m *model.AlgorithmVersionModel) *algorithm.Version {
	v := &algorithm.Version{
		ID:                    m.ID,
		AlgorithmID:           m.AlgorithmID,
		Version:               m.Version,
		Status:                algorithm.VersionStatus(m.Status),
		SelectionPolicy:       algorithm.SelectionPolicy(m.SelectionPolicy),
		DefaultImplementation: m.DefaultImplementationID,
		CreatedAt:             m.CreatedAt,
		UpdatedAt:             m.UpdatedAt,
	}
	for i := range m.Implementations {
		v.Implementations = append(v.Implementations, *AlgorithmImplementationToDomain(&m.Implementations[i]))
	}
	for i := range m.Evaluations {
		v.Evaluations = append(v.Evaluations, *AlgorithmEvaluationToDomain(&m.Evaluations[i]))
	}
	return v
}

func AlgorithmImplementationToModel(v *algorithm.Implementation) *model.AlgorithmImplementationModel {
	m := &model.AlgorithmImplementationModel{
		ID:           v.ID,
		VersionID:    v.VersionID,
		Name:         v.Name,
		Type:         string(v.Type),
		BindingRef:   v.BindingRef,
		LatencyMS:    v.LatencyMS,
		CostScore:    v.CostScore,
		QualityScore: v.QualityScore,
		Tier:         v.Tier,
		IsDefault:    v.IsDefault,
		CreatedAt:    v.CreatedAt,
		UpdatedAt:    v.UpdatedAt,
	}
	if v.Config != nil {
		b, _ := json.Marshal(v.Config)
		m.Config = datatypes.JSON(b)
	}
	return m
}

func AlgorithmImplementationToDomain(m *model.AlgorithmImplementationModel) *algorithm.Implementation {
	v := &algorithm.Implementation{
		ID:           m.ID,
		VersionID:    m.VersionID,
		Name:         m.Name,
		Type:         algorithm.ImplementationType(m.Type),
		BindingRef:   m.BindingRef,
		LatencyMS:    m.LatencyMS,
		CostScore:    m.CostScore,
		QualityScore: m.QualityScore,
		Tier:         m.Tier,
		IsDefault:    m.IsDefault,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
	if m.Config != nil {
		_ = json.Unmarshal(m.Config, &v.Config)
	}
	return v
}

func AlgorithmEvaluationToModel(v *algorithm.EvaluationProfile) *model.AlgorithmEvaluationModel {
	m := &model.AlgorithmEvaluationModel{
		ID:               v.ID,
		VersionID:        v.VersionID,
		DatasetRef:       v.DatasetRef,
		ReportArtifactID: v.ReportArtifactID,
		Summary:          v.Summary,
		CreatedAt:        v.CreatedAt,
		UpdatedAt:        v.UpdatedAt,
	}
	if v.Metrics != nil {
		b, _ := json.Marshal(v.Metrics)
		m.Metrics = datatypes.JSON(b)
	}
	return m
}

func AlgorithmEvaluationToDomain(m *model.AlgorithmEvaluationModel) *algorithm.EvaluationProfile {
	v := &algorithm.EvaluationProfile{
		ID:               m.ID,
		VersionID:        m.VersionID,
		DatasetRef:       m.DatasetRef,
		ReportArtifactID: m.ReportArtifactID,
		Summary:          m.Summary,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
	if m.Metrics != nil {
		_ = json.Unmarshal(m.Metrics, &v.Metrics)
	}
	return v
}
