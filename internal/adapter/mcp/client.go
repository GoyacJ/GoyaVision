package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"goyavision/internal/port"
)

// StaticClient 现用于承载真实 MCP(JSON-RPC) 协议适配实现。
// 为兼容既有注入链路，保留原类型名与构造函数。
type StaticClient struct {
	mu              sync.RWMutex
	servers         map[string]port.MCPServer
	meta            map[string]serverMeta
	configuredTools map[string][]port.MCPTool
	initialized     map[string]bool
	initMu          map[string]*sync.Mutex
	seq             int64
}

type serverMeta struct {
	Endpoint   string
	APIToken   string
	TimeoutSec int
}

type jsonRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int64       `json:"id,omitempty"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type jsonRPCResponse struct {
	JSONRPC string            `json:"jsonrpc"`
	ID      json.RawMessage   `json:"id"`
	Result  json.RawMessage   `json:"result"`
	Error   *jsonRPCErrorBody `json:"error"`
}

type jsonRPCErrorBody struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

var defaultClient = NewStaticClientWithoutDefaults()

func DefaultClient() *StaticClient {
	return defaultClient
}

func NewStaticClient() *StaticClient {
	// 为兼容旧调用保留该构造函数；默认不注入静态工具，统一走配置注册。
	return NewStaticClientWithoutDefaults()
}

func NewStaticClientWithoutDefaults() *StaticClient {
	return &StaticClient{
		servers:         make(map[string]port.MCPServer),
		meta:            make(map[string]serverMeta),
		configuredTools: make(map[string][]port.MCPTool),
		initialized:     make(map[string]bool),
		initMu:          make(map[string]*sync.Mutex),
	}
}

func (c *StaticClient) RegisterServer(server port.MCPServer, tools []port.MCPTool) {
	c.RegisterServerWithConfig(server, tools, "", "", 0)
}

func (c *StaticClient) RegisterServerWithConfig(server port.MCPServer, tools []port.MCPTool, endpoint, apiToken string, timeoutSec int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.servers[server.ID] = server
	c.configuredTools[server.ID] = append([]port.MCPTool(nil), tools...)
	c.meta[server.ID] = serverMeta{
		Endpoint:   strings.TrimSpace(endpoint),
		APIToken:   apiToken,
		TimeoutSec: timeoutSec,
	}
	if _, ok := c.initMu[server.ID]; !ok {
		c.initMu[server.ID] = &sync.Mutex{}
	}
	c.initialized[server.ID] = false
}

func (c *StaticClient) ListServers(_ context.Context) ([]port.MCPServer, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]port.MCPServer, 0, len(c.servers))
	for _, s := range c.servers {
		out = append(out, s)
	}
	return out, nil
}

func (c *StaticClient) GetServer(_ context.Context, serverID string) (*port.MCPServer, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	s, ok := c.servers[serverID]
	if !ok {
		return nil, fmt.Errorf("mcp server %s not found", serverID)
	}
	return &s, nil
}

func (c *StaticClient) ListTools(ctx context.Context, serverID string) ([]port.MCPTool, error) {
	meta, configured, err := c.getServerMeta(serverID)
	if err != nil {
		return nil, err
	}
	if meta.Endpoint == "" {
		// endpoint 为空时仅返回配置中的工具元信息（不可执行远程协议调用）
		out := make([]port.MCPTool, len(configured))
		copy(out, configured)
		return out, nil
	}

	if err := c.ensureInitialized(ctx, serverID, meta); err != nil {
		return nil, err
	}

	result, err := c.callRPC(ctx, meta, "tools/list", map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	tools, parseErr := parseMCPTools(result)
	if parseErr != nil {
		return nil, parseErr
	}
	return tools, nil
}

func (c *StaticClient) CallTool(ctx context.Context, serverID, toolName string, args map[string]interface{}) (map[string]interface{}, error) {
	meta, _, err := c.getServerMeta(serverID)
	if err != nil {
		return nil, err
	}
	if meta.Endpoint == "" {
		return nil, fmt.Errorf("mcp server %s endpoint is empty", serverID)
	}

	if err := c.ensureInitialized(ctx, serverID, meta); err != nil {
		return nil, err
	}

	if args == nil {
		args = map[string]interface{}{}
	}
	result, err := c.callRPC(ctx, meta, "tools/call", map[string]interface{}{
		"name":      toolName,
		"arguments": args,
	})
	if err != nil {
		return nil, err
	}

	out := map[string]interface{}{}
	if err := json.Unmarshal(result, &out); err != nil {
		return nil, fmt.Errorf("decode mcp tools/call result failed: %w", err)
	}
	if isErr, ok := out["isError"].(bool); ok && isErr {
		return nil, fmt.Errorf("mcp tool call returned error result: %v", out)
	}
	return out, nil
}

func (c *StaticClient) HealthCheck(ctx context.Context, serverID string) error {
	meta, _, err := c.getServerMeta(serverID)
	if err != nil {
		return err
	}
	if meta.Endpoint == "" {
		return fmt.Errorf("mcp server %s endpoint is empty", serverID)
	}
	return c.ensureInitialized(ctx, serverID, meta)
}

func (c *StaticClient) getServerMeta(serverID string) (serverMeta, []port.MCPTool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if _, ok := c.servers[serverID]; !ok {
		return serverMeta{}, nil, fmt.Errorf("mcp server %s not found", serverID)
	}
	meta := c.meta[serverID]
	tools := append([]port.MCPTool(nil), c.configuredTools[serverID]...)
	return meta, tools, nil
}

func (c *StaticClient) ensureInitialized(ctx context.Context, serverID string, meta serverMeta) error {
	c.mu.RLock()
	if c.initialized[serverID] {
		c.mu.RUnlock()
		return nil
	}
	mu := c.initMu[serverID]
	c.mu.RUnlock()

	if mu == nil {
		c.mu.Lock()
		if c.initMu[serverID] == nil {
			c.initMu[serverID] = &sync.Mutex{}
		}
		mu = c.initMu[serverID]
		c.mu.Unlock()
	}

	mu.Lock()
	defer mu.Unlock()

	c.mu.RLock()
	if c.initialized[serverID] {
		c.mu.RUnlock()
		return nil
	}
	c.mu.RUnlock()

	initParams := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]interface{}{},
		"clientInfo": map[string]interface{}{
			"name":    "goyavision",
			"version": "1.0.0",
		},
	}
	if _, err := c.callRPC(ctx, meta, "initialize", initParams); err != nil {
		return fmt.Errorf("mcp initialize failed: %w", err)
	}
	if err := c.notifyRPC(ctx, meta, "notifications/initialized", map[string]interface{}{}); err != nil {
		return fmt.Errorf("mcp initialized notification failed: %w", err)
	}

	c.mu.Lock()
	c.initialized[serverID] = true
	c.mu.Unlock()
	return nil
}

func (c *StaticClient) callRPC(ctx context.Context, meta serverMeta, method string, params interface{}) (json.RawMessage, error) {
	id := atomic.AddInt64(&c.seq, 1)
	reqBody := jsonRPCRequest{JSONRPC: "2.0", ID: id, Method: method, Params: params}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, meta.Endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if meta.APIToken != "" {
		req.Header.Set("Authorization", "Bearer "+meta.APIToken)
	}

	resp, err := buildHTTPClient(meta.TimeoutSec).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("mcp request failed: method=%s status=%d", method, resp.StatusCode)
	}

	var rpcResp jsonRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return nil, fmt.Errorf("decode mcp response failed: %w", err)
	}
	if rpcResp.Error != nil {
		return nil, fmt.Errorf("mcp rpc error: code=%d message=%s", rpcResp.Error.Code, rpcResp.Error.Message)
	}
	if len(rpcResp.Result) == 0 {
		return json.RawMessage(`{}`), nil
	}
	return rpcResp.Result, nil
}

func (c *StaticClient) notifyRPC(ctx context.Context, meta serverMeta, method string, params interface{}) error {
	reqBody := jsonRPCRequest{JSONRPC: "2.0", Method: method, Params: params}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, meta.Endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if meta.APIToken != "" {
		req.Header.Set("Authorization", "Bearer "+meta.APIToken)
	}

	resp, err := buildHTTPClient(meta.TimeoutSec).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("mcp notification failed: method=%s status=%d", method, resp.StatusCode)
	}
	return nil
}

func parseMCPTools(result json.RawMessage) ([]port.MCPTool, error) {
	var payload map[string]interface{}
	if err := json.Unmarshal(result, &payload); err != nil {
		return nil, fmt.Errorf("decode tools/list result failed: %w", err)
	}

	candidates := []interface{}{}
	if tools, ok := payload["tools"].([]interface{}); ok {
		candidates = tools
	} else if items, ok := payload["items"].([]interface{}); ok {
		candidates = items
	}

	out := make([]port.MCPTool, 0, len(candidates))
	for i := range candidates {
		m, ok := candidates[i].(map[string]interface{})
		if !ok {
			continue
		}
		tool := port.MCPTool{
			Name:        asString(m["name"]),
			Description: asString(m["description"]),
			Version:     asString(m["version"]),
		}
		if in, ok := m["inputSchema"].(map[string]interface{}); ok {
			tool.InputSchema = in
		} else if in, ok := m["input_schema"].(map[string]interface{}); ok {
			tool.InputSchema = in
		}
		if outSchema, ok := m["outputSchema"].(map[string]interface{}); ok {
			tool.OutputSchema = outSchema
		} else if outSchema, ok := m["output_schema"].(map[string]interface{}); ok {
			tool.OutputSchema = outSchema
		}
		if tool.Name != "" {
			out = append(out, tool)
		}
	}
	return out, nil
}

func buildHTTPClient(timeoutSec int) *http.Client {
	timeout := 15 * time.Second
	if timeoutSec > 0 {
		timeout = time.Duration(timeoutSec) * time.Second
	}
	return &http.Client{Timeout: timeout}
}

func joinURL(base, path string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(base))
	if err != nil {
		return "", fmt.Errorf("invalid mcp endpoint: %w", err)
	}
	u.Path = strings.TrimRight(u.Path, "/") + path
	return u.String(), nil
}

func asString(v interface{}) string {
	s, _ := v.(string)
	return s
}
