package mediamtx

import "time"

// PathConfig MediaMTX 路径配置
type PathConfig struct {
	Name   string `json:"name,omitempty"`
	Source string `json:"source,omitempty"`

	SourceFingerprint          string `json:"sourceFingerprint,omitempty"`
	SourceOnDemand             bool   `json:"sourceOnDemand,omitempty"`
	SourceOnDemandStartTimeout string `json:"sourceOnDemandStartTimeout,omitempty"`
	SourceOnDemandCloseAfter   string `json:"sourceOnDemandCloseAfter,omitempty"`
	MaxReaders                 int    `json:"maxReaders,omitempty"`

	RTSPTransport string `json:"rtspTransport,omitempty"`

	Record                *bool  `json:"record,omitempty"`
	RecordPath            string `json:"recordPath,omitempty"`
	RecordFormat          string `json:"recordFormat,omitempty"`
	RecordPartDuration    string `json:"recordPartDuration,omitempty"`
	RecordSegmentDuration string `json:"recordSegmentDuration,omitempty"`
	RecordDeleteAfter     string `json:"recordDeleteAfter,omitempty"`

	RunOnInit                  string `json:"runOnInit,omitempty"`
	RunOnInitRestart           bool   `json:"runOnInitRestart,omitempty"`
	RunOnDemand                string `json:"runOnDemand,omitempty"`
	RunOnDemandRestart         bool   `json:"runOnDemandRestart,omitempty"`
	RunOnReady                 string `json:"runOnReady,omitempty"`
	RunOnReadyRestart          bool   `json:"runOnReadyRestart,omitempty"`
	RunOnNotReady              string `json:"runOnNotReady,omitempty"`
	RunOnRead                  string `json:"runOnRead,omitempty"`
	RunOnReadRestart           bool   `json:"runOnReadRestart,omitempty"`
	RunOnUnread                string `json:"runOnUnread,omitempty"`
	RunOnRecordSegmentCreate   string `json:"runOnRecordSegmentCreate,omitempty"`
	RunOnRecordSegmentComplete string `json:"runOnRecordSegmentComplete,omitempty"`
}

// PathConfigList 路径配置列表响应
type PathConfigList struct {
	PageCount int          `json:"pageCount"`
	ItemCount int          `json:"itemCount"`
	Items     []PathConfig `json:"items"`
}

// Path 路径状态
type Path struct {
	Name          string       `json:"name"`
	ConfName      string       `json:"confName"`
	Source        *PathSource  `json:"source"`
	Ready         bool         `json:"ready"`
	ReadyTime     *time.Time   `json:"readyTime"`
	Available     bool         `json:"available"`
	AvailableTime *time.Time   `json:"availableTime"`
	Online        bool         `json:"online"`
	OnlineTime    *time.Time   `json:"onlineTime"`
	Tracks        []string     `json:"tracks"`
	BytesReceived uint64       `json:"bytesReceived"`
	BytesSent     uint64       `json:"bytesSent"`
	Readers       []PathReader `json:"readers"`
}

// PathSource 路径源信息
type PathSource struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// PathReader 路径读取者
type PathReader struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// PathList 路径列表响应
type PathList struct {
	PageCount int    `json:"pageCount"`
	ItemCount int    `json:"itemCount"`
	Items     []Path `json:"items"`
}

// Recording 录制信息
type Recording struct {
	Name     string             `json:"name"`
	Segments []RecordingSegment `json:"segments"`
}

// RecordingSegment 录制段信息
type RecordingSegment struct {
	Start time.Time `json:"start"`
}

// RecordingList 录制列表响应
type RecordingList struct {
	PageCount int         `json:"pageCount"`
	ItemCount int         `json:"itemCount"`
	Items     []Recording `json:"items"`
}

// HLSMuxer HLS 复用器信息
type HLSMuxer struct {
	Path        string    `json:"path"`
	Created     time.Time `json:"created"`
	LastRequest time.Time `json:"lastRequest"`
	BytesSent   uint64    `json:"bytesSent"`
}

// HLSMuxerList HLS 复用器列表响应
type HLSMuxerList struct {
	PageCount int        `json:"pageCount"`
	ItemCount int        `json:"itemCount"`
	Items     []HLSMuxer `json:"items"`
}

// RTSPSession RTSP 会话信息
type RTSPSession struct {
	ID                 string  `json:"id"`
	Created            string  `json:"created"`
	RemoteAddr         string  `json:"remoteAddr"`
	State              string  `json:"state"`
	Path               string  `json:"path"`
	Query              string  `json:"query"`
	Transport          string  `json:"transport"`
	BytesReceived      uint64  `json:"bytesReceived"`
	BytesSent          uint64  `json:"bytesSent"`
	RTPPacketsReceived uint64  `json:"rtpPacketsReceived"`
	RTPPacketsSent     uint64  `json:"rtpPacketsSent"`
	RTPPacketsLost     uint64  `json:"rtpPacketsLost"`
	RTPPacketsJitter   float64 `json:"rtpPacketsJitter"`
}

// RTSPSessionList RTSP 会话列表响应
type RTSPSessionList struct {
	PageCount int           `json:"pageCount"`
	ItemCount int           `json:"itemCount"`
	Items     []RTSPSession `json:"items"`
}

// RTMPConn RTMP 连接信息
type RTMPConn struct {
	ID            string `json:"id"`
	Created       string `json:"created"`
	RemoteAddr    string `json:"remoteAddr"`
	State         string `json:"state"`
	Path          string `json:"path"`
	Query         string `json:"query"`
	BytesReceived uint64 `json:"bytesReceived"`
	BytesSent     uint64 `json:"bytesSent"`
}

// RTMPConnList RTMP 连接列表响应
type RTMPConnList struct {
	PageCount int        `json:"pageCount"`
	ItemCount int        `json:"itemCount"`
	Items     []RTMPConn `json:"items"`
}

// Info MediaMTX 服务器信息
type Info struct {
	Version string `json:"version"`
	Started string `json:"started"`
}

// APIResponse 通用 API 响应
type APIResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}
