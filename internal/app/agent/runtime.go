package agentruntime

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"goyavision/internal/app/port"
	"goyavision/internal/domain/agent"
	"goyavision/internal/domain/workflow"

	"github.com/google/uuid"
)

type RunLoop struct {
	uow    port.UnitOfWork
	engine workflow.Engine
}

func NewRunLoop(uow port.UnitOfWork, workflowEngine workflow.Engine) *RunLoop {
	return &RunLoop{
		uow:    uow,
		engine: workflowEngine,
	}
}

func (r *RunLoop) StartSession(ctx context.Context, taskID uuid.UUID, budget map[string]interface{}) (*agent.Session, error) {
	now := time.Now().UTC()
	session := &agent.Session{
		ID:        uuid.New(),
		TaskID:    taskID,
		Status:    agent.SessionStatusRunning,
		Budget:    budget,
		StepCount: 0,
		StartedAt: now,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := r.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if repos.AgentSessions == nil {
			return nil
		}
		return repos.AgentSessions.Create(ctx, session)
	}); err != nil {
		return nil, err
	}
	return session, nil
}

func (r *RunLoop) FinishSession(ctx context.Context, sessionID uuid.UUID, status agent.SessionStatus) error {
	_, err := r.finishSession(ctx, sessionID, status)
	return err
}

func (r *RunLoop) RunStep(ctx context.Context, sessionID uuid.UUID, maxActions int) (*agent.Session, error) {
	if sessionID == uuid.Nil {
		return nil, errors.New("session_id is required")
	}
	if r.engine == nil {
		return nil, errors.New("workflow engine is not configured")
	}
	if maxActions <= 0 {
		maxActions = 1
	}

	var latestSession *agent.Session
	for i := 0; i < maxActions; i++ {
		session, task, wf, err := r.loadSessionTaskWorkflow(ctx, sessionID)
		if err != nil {
			return nil, err
		}
		latestSession = session

		if session.Status != agent.SessionStatusRunning {
			return session, nil
		}

		switch task.Status {
		case workflow.TaskStatusSuccess:
			_ = r.RecordEvent(ctx, &agent.RunEvent{
				TaskID:    task.ID,
				SessionID: sessionIDPtr(session.ID),
				EventType: agent.EventTypeAgentDecision,
				Source:    "agent_run_loop",
				Payload: map[string]interface{}{
					"decision":    "finish",
					"task_status": string(task.Status),
				},
			})
			latestSession, err = r.finishSession(ctx, session.ID, agent.SessionStatusSucceeded)
			if err != nil {
				return nil, err
			}
			return latestSession, nil
		case workflow.TaskStatusCancelled:
			_ = r.RecordEvent(ctx, &agent.RunEvent{
				TaskID:    task.ID,
				SessionID: sessionIDPtr(session.ID),
				EventType: agent.EventTypeAgentDecision,
				Source:    "agent_run_loop",
				Payload: map[string]interface{}{
					"decision":    "finish",
					"task_status": string(task.Status),
				},
			})
			latestSession, err = r.finishSession(ctx, session.ID, agent.SessionStatusCancelled)
			if err != nil {
				return nil, err
			}
			return latestSession, nil
		case workflow.TaskStatusPending, workflow.TaskStatusRunning:
			_ = r.RecordEvent(ctx, &agent.RunEvent{
				TaskID:    task.ID,
				SessionID: sessionIDPtr(session.ID),
				EventType: agent.EventTypeAgentDecision,
				Source:    "agent_run_loop",
				Payload: map[string]interface{}{
					"decision":    "wait",
					"task_status": string(task.Status),
				},
			})
			return latestSession, nil
		case workflow.TaskStatusFailed:
			toolErr := classifyTaskError(task.Error)
			policy := r.ResolveRecoveryPolicy(toolErr)
			_ = r.RecordEvent(ctx, &agent.RunEvent{
				TaskID:    task.ID,
				SessionID: sessionIDPtr(session.ID),
				EventType: agent.EventTypeAgentDecision,
				Source:    "agent_run_loop",
				Payload: map[string]interface{}{
					"decision":    "recover",
					"task_status": string(task.Status),
					"error": map[string]interface{}{
						"category":      string(toolErr.Category),
						"root_cause":    toolErr.RootCause,
						"action_hint":   toolErr.ActionHint,
						"retryable":     toolErr.Retryable,
						"provider_code": toolErr.ProviderCode,
						"message":       toolErr.Message,
					},
					"policy": map[string]interface{}{
						"retryable":     policy.Retryable,
						"backoff":       policy.Backoff,
						"fallback":      policy.Fallback,
						"require_human": policy.RequireHuman,
						"severity":      policy.Severity,
					},
				},
			})

			maxSteps := r.maxSteps(session.Budget)
			if !policy.Retryable || session.StepCount+1 > maxSteps {
				reason := "non_retryable"
				if session.StepCount+1 > maxSteps {
					reason = "budget_exhausted"
				}
				_ = r.RecordEvent(ctx, &agent.RunEvent{
					TaskID:    task.ID,
					SessionID: sessionIDPtr(session.ID),
					EventType: agent.EventTypeAgentEscalation,
					Source:    "agent_run_loop",
					Payload: map[string]interface{}{
						"reason":          reason,
						"category":        string(toolErr.Category),
						"step_count":      session.StepCount,
						"max_steps":       maxSteps,
						"require_human":   true,
						"recovery_policy": policy,
					},
				})
				latestSession, err = r.finishSession(ctx, session.ID, agent.SessionStatusFailed)
				if err != nil {
					return nil, err
				}
				return latestSession, nil
			}

			retryTask, updatedSession, err := r.createRetryTaskAndAdvanceSession(ctx, sessionID)
			if err != nil {
				return nil, err
			}
			latestSession = updatedSession

			_ = r.RecordEvent(ctx, &agent.RunEvent{
				TaskID:    retryTask.ID,
				SessionID: sessionIDPtr(updatedSession.ID),
				EventType: agent.EventTypeAgentAction,
				Source:    "agent_run_loop",
				Payload: map[string]interface{}{
					"action":         "retry_task",
					"from_task_id":   task.ID.String(),
					"retry_task_id":  retryTask.ID.String(),
					"attempt":        updatedSession.StepCount,
					"error_category": string(toolErr.Category),
					"fallback":       policy.Fallback,
				},
			})

			runErr := r.engine.Execute(ctx, wf, retryTask)
			if runErr != nil {
				_ = r.RecordEvent(ctx, &agent.RunEvent{
					TaskID:    retryTask.ID,
					SessionID: sessionIDPtr(updatedSession.ID),
					EventType: agent.EventTypeRecover,
					Source:    "agent_run_loop",
					Payload: map[string]interface{}{
						"result": "retry_failed",
						"error":  runErr.Error(),
					},
				})
				if markErr := r.markTaskFailedIfNeeded(ctx, retryTask.ID, runErr.Error()); markErr != nil {
					return nil, markErr
				}
			}

			sessionAfter, err := r.getSession(ctx, sessionID)
			if err == nil && sessionAfter != nil {
				latestSession = sessionAfter
			}

			switch retryTask.Status {
			case workflow.TaskStatusSuccess:
				_ = r.RecordEvent(ctx, &agent.RunEvent{
					TaskID:    retryTask.ID,
					SessionID: sessionIDPtr(updatedSession.ID),
					EventType: agent.EventTypeRecover,
					Source:    "agent_run_loop",
					Payload: map[string]interface{}{
						"result":        "recovered",
						"retry_task_id": retryTask.ID.String(),
					},
				})
				latestSession, err = r.finishSession(ctx, updatedSession.ID, agent.SessionStatusSucceeded)
				if err != nil {
					return nil, err
				}
				return latestSession, nil
			case workflow.TaskStatusCancelled:
				latestSession, err = r.finishSession(ctx, updatedSession.ID, agent.SessionStatusCancelled)
				if err != nil {
					return nil, err
				}
				return latestSession, nil
			default:
				// Keep session running. A next RunStep call can continue recovery attempts.
				continue
			}
		default:
			return latestSession, nil
		}
	}

	if latestSession == nil {
		return r.getSession(ctx, sessionID)
	}
	return latestSession, nil
}

func (r *RunLoop) RecordEvent(ctx context.Context, event *agent.RunEvent) error {
	if event == nil {
		return nil
	}
	if event.ID == uuid.Nil {
		event.ID = uuid.New()
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now().UTC()
	}
	return r.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if repos.RunEvents == nil {
			return nil
		}
		return repos.RunEvents.Create(ctx, event)
	})
}

func (r *RunLoop) ResolveRecoveryPolicy(toolErr *agent.ToolError) agent.RecoveryPolicy {
	if toolErr == nil {
		return agent.DefaultRecoveryPolicies[agent.ErrorCategoryUnknown]
	}
	if policy, ok := agent.DefaultRecoveryPolicies[toolErr.Category]; ok {
		if toolErr.Retryable {
			policy.Retryable = true
		}
		return policy
	}
	return agent.DefaultRecoveryPolicies[agent.ErrorCategoryUnknown]
}

func (r *RunLoop) loadSessionTaskWorkflow(ctx context.Context, sessionID uuid.UUID) (*agent.Session, *workflow.Task, *workflow.Workflow, error) {
	var (
		session *agent.Session
		task    *workflow.Task
		wf      *workflow.Workflow
	)
	err := r.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if repos.AgentSessions == nil {
			return errors.New("agent session repository is not configured")
		}
		if repos.Tasks == nil {
			return errors.New("task repository is not configured")
		}
		if repos.Workflows == nil {
			return errors.New("workflow repository is not configured")
		}

		var err error
		session, err = repos.AgentSessions.Get(ctx, sessionID)
		if err != nil {
			return err
		}
		task, err = repos.Tasks.Get(ctx, session.TaskID)
		if err != nil {
			return err
		}
		wf, err = repos.Workflows.GetWithNodes(ctx, task.WorkflowID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, nil, nil, err
	}
	return session, task, wf, nil
}

func (r *RunLoop) createRetryTaskAndAdvanceSession(ctx context.Context, sessionID uuid.UUID) (*workflow.Task, *agent.Session, error) {
	var (
		retryTask      *workflow.Task
		updatedSession *agent.Session
	)
	err := r.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if repos.AgentSessions == nil {
			return errors.New("agent session repository is not configured")
		}
		if repos.Tasks == nil {
			return errors.New("task repository is not configured")
		}

		session, err := repos.AgentSessions.Get(ctx, sessionID)
		if err != nil {
			return err
		}
		currentTask, err := repos.Tasks.Get(ctx, session.TaskID)
		if err != nil {
			return err
		}

		retryTask = &workflow.Task{
			ID:                 uuid.New(),
			WorkflowID:         currentTask.WorkflowID,
			WorkflowRevisionID: currentTask.WorkflowRevisionID,
			WorkflowRevision:   currentTask.WorkflowRevision,
			AssetID:            currentTask.AssetID,
			Status:             workflow.TaskStatusPending,
			Progress:           0,
			InputParams:        cloneMap(currentTask.InputParams),
		}
		if err := repos.Tasks.Create(ctx, retryTask); err != nil {
			return err
		}

		now := time.Now().UTC()
		session.StepCount++
		session.TaskID = retryTask.ID
		session.UpdatedAt = now
		if err := repos.AgentSessions.Update(ctx, session); err != nil {
			return err
		}
		updatedSession = session
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return retryTask, updatedSession, nil
}

func (r *RunLoop) finishSession(ctx context.Context, sessionID uuid.UUID, status agent.SessionStatus) (*agent.Session, error) {
	var result *agent.Session
	err := r.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if repos.AgentSessions == nil {
			return errors.New("agent session repository is not configured")
		}
		session, err := repos.AgentSessions.Get(ctx, sessionID)
		if err != nil {
			return err
		}
		now := time.Now().UTC()
		session.Status = status
		session.EndedAt = &now
		session.UpdatedAt = now
		if err := repos.AgentSessions.Update(ctx, session); err != nil {
			return err
		}
		result = session
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *RunLoop) getSession(ctx context.Context, sessionID uuid.UUID) (*agent.Session, error) {
	var session *agent.Session
	err := r.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if repos.AgentSessions == nil {
			return errors.New("agent session repository is not configured")
		}
		var err error
		session, err = repos.AgentSessions.Get(ctx, sessionID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (r *RunLoop) markTaskFailedIfNeeded(ctx context.Context, taskID uuid.UUID, errMsg string) error {
	return r.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if repos.Tasks == nil {
			return errors.New("task repository is not configured")
		}
		task, err := repos.Tasks.Get(ctx, taskID)
		if err != nil {
			return err
		}
		if task.IsCompleted() {
			return nil
		}
		now := time.Now().UTC()
		task.Status = workflow.TaskStatusFailed
		task.Error = errMsg
		task.CompletedAt = &now
		return repos.Tasks.Update(ctx, task)
	})
}

func (r *RunLoop) maxSteps(budget map[string]interface{}) int {
	const defaultMaxSteps = 3
	if budget == nil {
		return defaultMaxSteps
	}
	if v, ok := budget["max_steps"]; ok {
		if parsed := numberToInt(v); parsed > 0 {
			return parsed
		}
	}
	if v, ok := budget["max_retries"]; ok {
		if parsed := numberToInt(v); parsed > 0 {
			return parsed
		}
	}
	return defaultMaxSteps
}

func numberToInt(v interface{}) int {
	switch n := v.(type) {
	case int:
		return n
	case int32:
		return int(n)
	case int64:
		return int(n)
	case float32:
		return int(math.Round(float64(n)))
	case float64:
		return int(math.Round(n))
	case string:
		parsed, err := strconv.Atoi(strings.TrimSpace(n))
		if err == nil {
			return parsed
		}
	}
	return 0
}

func cloneMap(src map[string]interface{}) map[string]interface{} {
	if src == nil {
		return nil
	}
	dst := make(map[string]interface{}, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func classifyTaskError(message string) *agent.ToolError {
	msg := strings.TrimSpace(message)
	lower := strings.ToLower(msg)
	toolErr := &agent.ToolError{
		Category:   agent.ErrorCategoryUnknown,
		RootCause:  "dependency",
		ActionHint: "collect_diagnostics",
		Message:    msg,
	}

	switch {
	case containsAny(lower, "policy denied", "permission denied", "forbidden", "unauthorized"):
		toolErr.Category = agent.ErrorCategoryPolicyDeny
		toolErr.RootCause = "permission"
		toolErr.ActionHint = "request_approval_or_scope_update"
	case containsAny(lower, "deadline exceeded", "timeout"):
		toolErr.Category = agent.ErrorCategoryTimeout
		toolErr.RootCause = "dependency"
		toolErr.ActionHint = "increase_timeout_or_reduce_payload"
	case containsAny(lower, "invalid", "validation", "schema"):
		toolErr.Category = agent.ErrorCategoryValidation
		toolErr.RootCause = "input"
		toolErr.ActionHint = "repair_input_or_mapping"
	case containsAny(lower, "quota", "rate limit", "resource exhausted", "out of memory"):
		toolErr.Category = agent.ErrorCategoryResourceLimit
		toolErr.RootCause = "resource"
		toolErr.ActionHint = "throttle_or_switch_low_cost_impl"
	case containsAny(lower, "context version conflict", "temporarily unavailable", "temporary"):
		toolErr.Category = agent.ErrorCategoryTransient
		toolErr.RootCause = "dependency"
		toolErr.ActionHint = "retry_with_backoff"
	case containsAny(lower, "connection refused", "connection reset", "dns", "i/o timeout", "service unavailable", "503"):
		toolErr.Category = agent.ErrorCategoryDependency
		toolErr.RootCause = "dependency"
		toolErr.ActionHint = "retry_or_switch_dependency"
	case containsAny(lower, "panic", "segmentation fault", "nil pointer"):
		toolErr.Category = agent.ErrorCategoryToolBug
		toolErr.RootCause = "dependency"
		toolErr.ActionHint = "switch_impl_and_report_bug"
	case containsAny(lower, "reasoning", "hallucination", "model output invalid"):
		toolErr.Category = agent.ErrorCategoryModelReasoning
		toolErr.RootCause = "input"
		toolErr.ActionHint = "switch_model_or_prompt"
	}

	if policy, ok := agent.DefaultRecoveryPolicies[toolErr.Category]; ok {
		toolErr.Retryable = policy.Retryable
	}
	if toolErr.Message == "" {
		toolErr.Message = fmt.Sprintf("task failed with category %s", toolErr.Category)
	}
	return toolErr
}

func containsAny(s string, terms ...string) bool {
	for _, term := range terms {
		if strings.Contains(s, term) {
			return true
		}
	}
	return false
}

func sessionIDPtr(id uuid.UUID) *uuid.UUID {
	out := id
	return &out
}
