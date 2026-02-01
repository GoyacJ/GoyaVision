package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   Server
	DB       DB
	FFmpeg   FFmpeg
	Preview  Preview
	Record   Record
	AI       AI
	JWT      JWT
	MediaMTX MediaMTX
}

type MediaMTX struct {
	APIAddress      string
	RTSPAddress     string
	RTMPAddress     string
	HLSAddress      string
	WebRTCAddress   string
	PlaybackAddress string
	RecordPath      string
	RecordFormat    string
	SegmentDuration string
}

type JWT struct {
	Secret     string
	Expire     time.Duration
	RefreshExp time.Duration
	Issuer     string
}

type Server struct {
	Port int
}

type DB struct {
	DSN string
}

type FFmpeg struct {
	Bin       string
	MaxRecord int
	MaxFrame  int
}

type Preview struct {
	Provider    string
	MediamtxBin string
	MaxPreview  int
	HLSBase     string
}

type Record struct {
	BasePath   string
	SegmentSec int
}

type AI struct {
	Timeout time.Duration
	Retry   int
}

func (s Server) Addr() string {
	return fmt.Sprintf(":%d", s.Port)
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	v.SetEnvPrefix("GOYAVISION")
	v.AutomaticEnv()

	cfg := &Config{
		Server: Server{Port: v.GetInt("server.port")},
		DB:     DB{DSN: v.GetString("db.dsn")},
		FFmpeg: FFmpeg{
			Bin:       v.GetString("ffmpeg.bin"),
			MaxRecord: v.GetInt("ffmpeg.max_record"),
			MaxFrame:  v.GetInt("ffmpeg.max_frame"),
		},
		Preview: Preview{
			Provider:    v.GetString("preview.provider"),
			MediamtxBin: v.GetString("preview.mediamtx_bin"),
			MaxPreview:  v.GetInt("preview.max_preview"),
			HLSBase:     v.GetString("preview.hls_base"),
		},
		Record: Record{
			BasePath:   v.GetString("record.base_path"),
			SegmentSec: v.GetInt("record.segment_sec"),
		},
		AI: AI{
			Timeout: v.GetDuration("ai.timeout"),
			Retry:   v.GetInt("ai.retry"),
		},
		JWT: JWT{
			Secret:     v.GetString("jwt.secret"),
			Expire:     v.GetDuration("jwt.expire"),
			RefreshExp: v.GetDuration("jwt.refresh_exp"),
			Issuer:     v.GetString("jwt.issuer"),
		},
		MediaMTX: MediaMTX{
			APIAddress:      v.GetString("mediamtx.api_address"),
			RTSPAddress:     v.GetString("mediamtx.rtsp_address"),
			RTMPAddress:     v.GetString("mediamtx.rtmp_address"),
			HLSAddress:      v.GetString("mediamtx.hls_address"),
			WebRTCAddress:   v.GetString("mediamtx.webrtc_address"),
			PlaybackAddress: v.GetString("mediamtx.playback_address"),
			RecordPath:      v.GetString("mediamtx.record_path"),
			RecordFormat:    v.GetString("mediamtx.record_format"),
			SegmentDuration: v.GetString("mediamtx.segment_duration"),
		},
	}
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.JWT.Secret == "" {
		cfg.JWT.Secret = "goyavision-secret-key"
	}
	if cfg.JWT.Expire == 0 {
		cfg.JWT.Expire = 2 * time.Hour
	}
	if cfg.JWT.RefreshExp == 0 {
		cfg.JWT.RefreshExp = 7 * 24 * time.Hour
	}
	if cfg.JWT.Issuer == "" {
		cfg.JWT.Issuer = "goyavision"
	}
	if cfg.MediaMTX.APIAddress == "" {
		cfg.MediaMTX.APIAddress = "http://localhost:9997"
	}
	if cfg.MediaMTX.RTSPAddress == "" {
		cfg.MediaMTX.RTSPAddress = "rtsp://localhost:8554"
	}
	if cfg.MediaMTX.RTMPAddress == "" {
		cfg.MediaMTX.RTMPAddress = "rtmp://localhost:1935"
	}
	if cfg.MediaMTX.HLSAddress == "" {
		cfg.MediaMTX.HLSAddress = "http://localhost:8888"
	}
	if cfg.MediaMTX.WebRTCAddress == "" {
		cfg.MediaMTX.WebRTCAddress = "http://localhost:8889"
	}
	if cfg.MediaMTX.PlaybackAddress == "" {
		cfg.MediaMTX.PlaybackAddress = "http://localhost:9996"
	}
	if cfg.MediaMTX.RecordPath == "" {
		cfg.MediaMTX.RecordPath = "./data/recordings/%path/%Y-%m-%d_%H-%M-%S"
	}
	if cfg.MediaMTX.RecordFormat == "" {
		cfg.MediaMTX.RecordFormat = "fmp4"
	}
	if cfg.MediaMTX.SegmentDuration == "" {
		cfg.MediaMTX.SegmentDuration = "1h"
	}
	return cfg, nil
}
