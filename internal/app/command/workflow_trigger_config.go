package command

import (
	"fmt"

	"goyavision/internal/domain/workflow"
)

func buildTriggerConfig(raw map[string]interface{}) (*workflow.TriggerConfig, error) {
	if raw == nil {
		return nil, nil
	}

	tc := &workflow.TriggerConfig{}

	if v, ok := raw["schedule"]; ok {
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("trigger_conf.schedule must be string")
		}
		tc.Schedule = s
	}

	if v, ok := raw["interval_sec"]; ok {
		switch iv := v.(type) {
		case float64:
			tc.IntervalSec = int(iv)
		case int:
			tc.IntervalSec = iv
		default:
			return nil, fmt.Errorf("trigger_conf.interval_sec must be number")
		}
	}

	if v, ok := raw["event_type"]; ok {
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("trigger_conf.event_type must be string")
		}
		tc.EventType = s
	}

	if v, ok := raw["event_filter"]; ok {
		m, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("trigger_conf.event_filter must be object")
		}
		tc.EventFilter = m
	}

	return tc, nil
}
