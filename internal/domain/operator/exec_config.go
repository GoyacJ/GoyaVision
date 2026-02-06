package operator

type ExecConfig struct {
	HTTP *HTTPExecConfig `json:"http,omitempty"`
	CLI  *CLIExecConfig  `json:"cli,omitempty"`
	MCP  *MCPExecConfig  `json:"mcp,omitempty"`
}

type HTTPExecConfig struct {
	Endpoint   string            `json:"endpoint"`
	Method     string            `json:"method"`
	Headers    map[string]string `json:"headers,omitempty"`
	TimeoutSec int               `json:"timeout_sec,omitempty"`
	AuthType   string            `json:"auth_type,omitempty"`
	AuthConfig map[string]string `json:"auth_config,omitempty"`
}

type CLIExecConfig struct {
	Command    string            `json:"command"`
	Args       []string          `json:"args"`
	WorkDir    string            `json:"work_dir,omitempty"`
	Env        map[string]string `json:"env,omitempty"`
	TimeoutSec int               `json:"timeout_sec,omitempty"`
}

type MCPExecConfig struct {
	ServerID      string                 `json:"server_id"`
	ToolName      string                 `json:"tool_name"`
	ToolVersion   string                 `json:"tool_version,omitempty"`
	TimeoutSec    int                    `json:"timeout_sec,omitempty"`
	InputMapping  map[string]interface{} `json:"input_mapping,omitempty"`
	OutputMapping map[string]interface{} `json:"output_mapping,omitempty"`
}
