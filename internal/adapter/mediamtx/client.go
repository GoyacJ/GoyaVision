package mediamtx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client MediaMTX HTTP API 客户端
type Client struct {
	baseURL    string
	username   string
	password   string
	httpClient *http.Client
}

// NewClient 创建 MediaMTX 客户端
func NewClient(baseURL, username, password string) *Client {
	return &Client{
		baseURL:  baseURL,
		username: username,
		password: password,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetInfo 获取服务器信息
func (c *Client) GetInfo(ctx context.Context) (*Info, error) {
	var info Info
	err := c.doRequest(ctx, http.MethodGet, "/v3/info", nil, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// AddPath 添加路径配置
func (c *Client) AddPath(ctx context.Context, name string, cfg *PathConfig) error {
	path := fmt.Sprintf("/v3/config/paths/add/%s", url.PathEscape(name))
	var resp APIResponse
	err := c.doRequest(ctx, http.MethodPost, path, cfg, &resp)
	if err != nil {
		return err
	}
	if resp.Status == "error" {
		return fmt.Errorf("mediamtx error: %s", resp.Error)
	}
	return nil
}

// GetPathConfig 获取路径配置
func (c *Client) GetPathConfig(ctx context.Context, name string) (*PathConfig, error) {
	path := fmt.Sprintf("/v3/config/paths/get/%s", url.PathEscape(name))
	var cfg PathConfig
	err := c.doRequest(ctx, http.MethodGet, path, nil, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// PatchPath 更新路径配置
func (c *Client) PatchPath(ctx context.Context, name string, cfg *PathConfig) error {
	path := fmt.Sprintf("/v3/config/paths/patch/%s", url.PathEscape(name))
	var resp APIResponse
	err := c.doRequest(ctx, http.MethodPatch, path, cfg, &resp)
	if err != nil {
		return err
	}
	if resp.Status == "error" {
		return fmt.Errorf("mediamtx error: %s", resp.Error)
	}
	return nil
}

// DeletePath 删除路径配置
func (c *Client) DeletePath(ctx context.Context, name string) error {
	path := fmt.Sprintf("/v3/config/paths/delete/%s", url.PathEscape(name))
	var resp APIResponse
	err := c.doRequest(ctx, http.MethodDelete, path, nil, &resp)
	if err != nil {
		return err
	}
	if resp.Status == "error" {
		return fmt.Errorf("mediamtx error: %s", resp.Error)
	}
	return nil
}

// ListPathConfigs 列出所有路径配置
func (c *Client) ListPathConfigs(ctx context.Context, page, itemsPerPage int) (*PathConfigList, error) {
	path := fmt.Sprintf("/v3/config/paths/list?page=%d&itemsPerPage=%d", page, itemsPerPage)
	var list PathConfigList
	err := c.doRequest(ctx, http.MethodGet, path, nil, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// GetPath 获取路径状态
func (c *Client) GetPath(ctx context.Context, name string) (*Path, error) {
	path := fmt.Sprintf("/v3/paths/get/%s", url.PathEscape(name))
	var p Path
	err := c.doRequest(ctx, http.MethodGet, path, nil, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// ListPaths 列出所有路径状态
func (c *Client) ListPaths(ctx context.Context, page, itemsPerPage int) (*PathList, error) {
	path := fmt.Sprintf("/v3/paths/list?page=%d&itemsPerPage=%d", page, itemsPerPage)
	var list PathList
	err := c.doRequest(ctx, http.MethodGet, path, nil, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// IsPathReady 检查路径是否就绪
func (c *Client) IsPathReady(ctx context.Context, name string) (bool, error) {
	p, err := c.GetPath(ctx, name)
	if err != nil {
		return false, err
	}
	return p.Ready, nil
}

// GetRecordings 获取路径的录制列表
func (c *Client) GetRecordings(ctx context.Context, name string) (*Recording, error) {
	path := fmt.Sprintf("/v3/recordings/get/%s", url.PathEscape(name))
	var rec Recording
	err := c.doRequest(ctx, http.MethodGet, path, nil, &rec)
	if err != nil {
		return nil, err
	}
	return &rec, nil
}

// ListRecordings 列出所有录制
func (c *Client) ListRecordings(ctx context.Context, page, itemsPerPage int) (*RecordingList, error) {
	path := fmt.Sprintf("/v3/recordings/list?page=%d&itemsPerPage=%d", page, itemsPerPage)
	var list RecordingList
	err := c.doRequest(ctx, http.MethodGet, path, nil, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// DeleteRecordingSegment 删除录制段
func (c *Client) DeleteRecordingSegment(ctx context.Context, pathName string, start time.Time) error {
	path := fmt.Sprintf("/v3/recordings/deletesegment?path=%s&start=%s",
		url.QueryEscape(pathName), url.QueryEscape(start.Format(time.RFC3339)))
	var resp APIResponse
	err := c.doRequest(ctx, http.MethodDelete, path, nil, &resp)
	if err != nil {
		return err
	}
	if resp.Status == "error" {
		return fmt.Errorf("mediamtx error: %s", resp.Error)
	}
	return nil
}

// ListHLSMuxers 列出 HLS 复用器
func (c *Client) ListHLSMuxers(ctx context.Context, page, itemsPerPage int) (*HLSMuxerList, error) {
	path := fmt.Sprintf("/v3/hlsmuxers/list?page=%d&itemsPerPage=%d", page, itemsPerPage)
	var list HLSMuxerList
	err := c.doRequest(ctx, http.MethodGet, path, nil, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// GetHLSMuxer 获取 HLS 复用器信息
func (c *Client) GetHLSMuxer(ctx context.Context, name string) (*HLSMuxer, error) {
	path := fmt.Sprintf("/v3/hlsmuxers/get/%s", url.PathEscape(name))
	var muxer HLSMuxer
	err := c.doRequest(ctx, http.MethodGet, path, nil, &muxer)
	if err != nil {
		return nil, err
	}
	return &muxer, nil
}

// ListRTSPSessions 列出 RTSP 会话
func (c *Client) ListRTSPSessions(ctx context.Context, page, itemsPerPage int) (*RTSPSessionList, error) {
	path := fmt.Sprintf("/v3/rtspsessions/list?page=%d&itemsPerPage=%d", page, itemsPerPage)
	var list RTSPSessionList
	err := c.doRequest(ctx, http.MethodGet, path, nil, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// KickRTSPSession 踢出 RTSP 会话
func (c *Client) KickRTSPSession(ctx context.Context, id string) error {
	path := fmt.Sprintf("/v3/rtspsessions/kick/%s", url.PathEscape(id))
	var resp APIResponse
	err := c.doRequest(ctx, http.MethodPost, path, nil, &resp)
	if err != nil {
		return err
	}
	if resp.Status == "error" {
		return fmt.Errorf("mediamtx error: %s", resp.Error)
	}
	return nil
}

// ListRTMPConns 列出 RTMP 连接
func (c *Client) ListRTMPConns(ctx context.Context, page, itemsPerPage int) (*RTMPConnList, error) {
	path := fmt.Sprintf("/v3/rtmpconns/list?page=%d&itemsPerPage=%d", page, itemsPerPage)
	var list RTMPConnList
	err := c.doRequest(ctx, http.MethodGet, path, nil, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// KickRTMPConn 踢出 RTMP 连接
func (c *Client) KickRTMPConn(ctx context.Context, id string) error {
	path := fmt.Sprintf("/v3/rtmpconns/kick/%s", url.PathEscape(id))
	var resp APIResponse
	err := c.doRequest(ctx, http.MethodPost, path, nil, &resp)
	if err != nil {
		return err
	}
	if resp.Status == "error" {
		return fmt.Errorf("mediamtx error: %s", resp.Error)
	}
	return nil
}

// EnableRecording 启用路径录制
func (c *Client) EnableRecording(ctx context.Context, name string, recordPath, format, segmentDuration string) error {
	record := true
	cfg := &PathConfig{
		Record:                &record,
		RecordPath:            recordPath,
		RecordFormat:          format,
		RecordSegmentDuration: segmentDuration,
	}
	return c.PatchPath(ctx, name, cfg)
}

// DisableRecording 禁用路径录制
func (c *Client) DisableRecording(ctx context.Context, name string) error {
	record := false
	cfg := &PathConfig{
		Record: &record,
	}
	return c.PatchPath(ctx, name, cfg)
}

// doRequest 执行 HTTP 请求
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	fullURL := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.username != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var apiResp APIResponse
		if err := json.Unmarshal(respBody, &apiResp); err == nil && apiResp.Error != "" {
			return fmt.Errorf("mediamtx error (status %d): %s", resp.StatusCode, apiResp.Error)
		}
		return fmt.Errorf("mediamtx error: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("unmarshal response: %w", err)
		}
	}

	return nil
}

// Ping 检查 MediaMTX 服务是否可用
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.GetInfo(ctx)
	return err
}
