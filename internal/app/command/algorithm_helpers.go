package command

import (
	"goyavision/internal/app/dto"
	"goyavision/internal/domain/algorithm"
	"goyavision/pkg/apperr"

	"github.com/google/uuid"
)

func buildAlgorithmVersionFromInput(
	input *dto.AlgorithmVersionInput,
	algorithmID uuid.UUID,
) (*algorithm.Version, []algorithm.Implementation, []algorithm.EvaluationProfile, error) {
	if input == nil {
		return nil, nil, nil, apperr.InvalidInput("initial_version is required")
	}
	version := input.Version
	if version == "" {
		return nil, nil, nil, apperr.InvalidInput("version is required")
	}
	status := input.Status
	if status == "" {
		status = algorithm.VersionStatusDraft
	}
	selection := input.SelectionPolicy
	if selection == "" {
		selection = algorithm.SelectionPolicyStable
	}

	v := &algorithm.Version{
		ID:              uuid.New(),
		AlgorithmID:     algorithmID,
		Version:         version,
		Status:          status,
		SelectionPolicy: selection,
	}

	impls, defaultID, err := mapAlgorithmImplementations(input.Implementations, v.ID)
	if err != nil {
		return nil, nil, nil, err
	}
	v.DefaultImplementation = defaultID

	evals := mapAlgorithmEvaluations(input.Evaluations, v.ID)
	return v, impls, evals, nil
}

func buildAlgorithmVersionFromCommand(
	cmd dto.CreateAlgorithmVersionCommand,
) (*algorithm.Version, []algorithm.Implementation, []algorithm.EvaluationProfile, error) {
	if cmd.Version == "" {
		return nil, nil, nil, apperr.InvalidInput("version is required")
	}
	status := cmd.Status
	if status == "" {
		status = algorithm.VersionStatusDraft
	}
	selection := cmd.SelectionPolicy
	if selection == "" {
		selection = algorithm.SelectionPolicyStable
	}

	v := &algorithm.Version{
		ID:              uuid.New(),
		AlgorithmID:     cmd.AlgorithmID,
		Version:         cmd.Version,
		Status:          status,
		SelectionPolicy: selection,
	}

	impls, defaultID, err := mapAlgorithmImplementations(cmd.Implementations, v.ID)
	if err != nil {
		return nil, nil, nil, err
	}
	v.DefaultImplementation = defaultID

	evals := mapAlgorithmEvaluations(cmd.Evaluations, v.ID)
	return v, impls, evals, nil
}

func mapAlgorithmImplementations(inputs []dto.AlgorithmImplementationInput, versionID uuid.UUID) ([]algorithm.Implementation, *uuid.UUID, error) {
	if len(inputs) == 0 {
		return nil, nil, apperr.InvalidInput("at least one implementation is required")
	}

	impls := make([]algorithm.Implementation, 0, len(inputs))
	var defaultID *uuid.UUID
	for i := range inputs {
		if inputs[i].BindingRef == "" {
			return nil, nil, apperr.InvalidInput("implementation.binding_ref is required")
		}

		implType := inputs[i].Type
		if implType == "" {
			implType = algorithm.ImplementationOperatorVersion
		}
		impl := algorithm.Implementation{
			ID:           uuid.New(),
			VersionID:    versionID,
			Name:         inputs[i].Name,
			Type:         implType,
			BindingRef:   inputs[i].BindingRef,
			Config:       inputs[i].Config,
			LatencyMS:    inputs[i].LatencyMS,
			CostScore:    inputs[i].CostScore,
			QualityScore: inputs[i].QualityScore,
			Tier:         inputs[i].Tier,
			IsDefault:    inputs[i].IsDefault,
		}
		if impl.Tier == "" {
			impl.Tier = string(algorithm.SelectionPolicyStable)
		}

		if impl.IsDefault {
			if defaultID != nil {
				return nil, nil, apperr.InvalidInput("multiple default implementations are not allowed")
			}
			id := impl.ID
			defaultID = &id
		}
		impls = append(impls, impl)
	}

	if defaultID == nil && len(impls) > 0 {
		impls[0].IsDefault = true
		id := impls[0].ID
		defaultID = &id
	}
	return impls, defaultID, nil
}

func mapAlgorithmEvaluations(inputs []dto.AlgorithmEvaluationInput, versionID uuid.UUID) []algorithm.EvaluationProfile {
	if len(inputs) == 0 {
		return nil
	}
	out := make([]algorithm.EvaluationProfile, 0, len(inputs))
	for i := range inputs {
		out = append(out, algorithm.EvaluationProfile{
			ID:               uuid.New(),
			VersionID:        versionID,
			DatasetRef:       inputs[i].DatasetRef,
			Metrics:          inputs[i].Metrics,
			ReportArtifactID: inputs[i].ReportArtifactID,
			Summary:          inputs[i].Summary,
		})
	}
	return out
}
