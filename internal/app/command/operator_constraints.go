package command

import (
	"regexp"

	"goyavision/internal/domain/operator"
	"goyavision/pkg/apperr"
)

var semverPattern = regexp.MustCompile(`^v?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-[0-9A-Za-z.-]+)?(?:\+[0-9A-Za-z.-]+)?$`)

func validateExecMode(mode operator.ExecMode) error {
	if mode == operator.ExecModeHTTP || mode == operator.ExecModeCLI || mode == operator.ExecModeMCP {
		return nil
	}
	return apperr.InvalidInput("invalid exec_mode, allowed values: http|cli|mcp")
}

func validateSemver(version string) error {
	if !semverPattern.MatchString(version) {
		return apperr.InvalidInput("invalid version, must be semver like 1.0.0")
	}
	return nil
}
