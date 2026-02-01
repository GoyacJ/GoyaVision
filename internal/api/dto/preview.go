package dto

type PreviewStartResponse struct {
	HLSURL    string `json:"hls_url"`
	RTSPUrl   string `json:"rtsp_url"`
	RTMPUrl   string `json:"rtmp_url"`
	WebRTCUrl string `json:"webrtc_url"`
}
