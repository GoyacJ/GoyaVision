package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"goyavision/internal/domain/operator"
	"goyavision/internal/port"
)

var _ port.OperatorExecutor = (*CLIOperatorExecutor)(nil)

// CLIOperatorExecutor CLI 算子执行器
type CLIOperatorExecutor struct{}

func NewCLIOperatorExecutor() *CLIOperatorExecutor {
	return &CLIOperatorExecutor{}
}

func (e *CLIOperatorExecutor) Execute(ctx context.Context, version *operator.OperatorVersion, input *operator.Input) (*operator.Output, error) {
	if version == nil {
		return nil, fmt.Errorf("operator version is nil")
	}
	if version.ExecMode != operator.ExecModeCLI {
		return nil, fmt.Errorf("cli executor does not support exec mode: %s", version.ExecMode)
	}
	if version.ExecConfig == nil || version.ExecConfig.CLI == nil {
		return nil, fmt.Errorf("cli exec config is required")
	}

	cliCfg := version.ExecConfig.CLI
	if strings.TrimSpace(cliCfg.Command) == "" {
		return nil, fmt.Errorf("cli command is required")
	}

	execCtx := ctx
	if cliCfg.TimeoutSec > 0 {
		var cancel context.CancelFunc
		execCtx, cancel = context.WithTimeout(ctx, time.Duration(cliCfg.TimeoutSec)*time.Second)
		defer cancel()
	}

	cmd := exec.CommandContext(execCtx, cliCfg.Command, cliCfg.Args...)
	if strings.TrimSpace(cliCfg.WorkDir) != "" {
		cmd.Dir = cliCfg.WorkDir
	}

	if len(cliCfg.Env) > 0 {
		env := os.Environ()
		for k, v := range cliCfg.Env {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
		cmd.Env = env
	}

	if input == nil {
		input = &operator.Input{}
	}
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	cmd.Stdin = bytes.NewReader(inputBytes)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("cli command failed: %w, stderr: %s", err, strings.TrimSpace(stderr.String()))
	}

	out := strings.TrimSpace(stdout.String())
	if out == "" {
		return &operator.Output{}, nil
	}

	var output operator.Output
	if err := json.Unmarshal([]byte(out), &output); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cli output: %w", err)
	}

	return &output, nil
}

func (e *CLIOperatorExecutor) Mode() operator.ExecMode {
	return operator.ExecModeCLI
}

func (e *CLIOperatorExecutor) HealthCheck(ctx context.Context, version *operator.OperatorVersion) error {
	if version == nil {
		return fmt.Errorf("operator version is nil")
	}
	if version.ExecConfig == nil || version.ExecConfig.CLI == nil {
		return fmt.Errorf("cli exec config is required")
	}
	if strings.TrimSpace(version.ExecConfig.CLI.Command) == "" {
		return fmt.Errorf("cli command is required")
	}
	return nil
}
