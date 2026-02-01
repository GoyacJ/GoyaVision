package dto

import "time"

type PlaybackURLsResponse struct {
	HLS string `json:"hls_url"`
	MP4 string `json:"mp4_url"`
}

type PlaybackSegmentResponse struct {
	Start       time.Time `json:"start"`
	PlaybackURL string    `json:"playback_url"`
}

type PlaybackRequest struct {
	Start string `query:"start"`
}
