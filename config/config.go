package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Env        string
	Server     Server
	DB         DB
	FFmpeg     FFmpeg
	Preview    Preview
	Record     Record
	AI         AI
	JWT        JWT
	MediaMTX   MediaMTX
	MinIO      MinIO
	MCP        MCP
	OAuth      OAuth
	Payment    Payment
	EncryptKey string
}

type Payment struct {
	Alipay AlipayConfig `mapstructure:"alipay"`
	Wechat WechatConfig `mapstructure:"wechat"`
	Union  UnionConfig  `mapstructure:"union"`
}

type AlipayConfig struct {
	AppID      string `mapstructure:"app_id"`
	PrivateKey string `mapstructure:"private_key"`
	PublicKey  string `mapstructure:"public_key"`
	NotifyURL  string `mapstructure:"notify_url"`
	ReturnURL  string `mapstructure:"return_url"`
	IsProd     bool   `mapstructure:"is_prod"`
}

type WechatConfig struct {
	AppID     string `mapstructure:"app_id"`
	MchID     string `mapstructure:"mch_id"`
	APIKey    string `mapstructure:"api_key"`
	CertPath  string `mapstructure:"cert_path"`
	KeyPath   string `mapstructure:"key_path"`
	NotifyURL string `mapstructure:"notify_url"`
	IsProd    bool   `mapstructure:"is_prod"`
}

type UnionConfig struct {
	MerID     string `mapstructure:"mer_id"`
	CertPath  string `mapstructure:"cert_path"`
	NotifyURL string `mapstructure:"notify_url"`
	IsProd    bool   `mapstructure:"is_prod"`
}

type OAuth struct {
	Github OAuthConfig `mapstructure:"github"`
	Wechat OAuthConfig `mapstructure:"wechat"`
}

type OAuthConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURI  string `mapstructure:"redirect_uri"`
}

type MCP struct {
	Servers []MCPServer `mapstructure:"servers"`
}

type MCPServer struct {
	ID          string    `mapstructure:"id"`
	Name        string    `mapstructure:"name"`
	Description string    `mapstructure:"description"`
	Status      string    `mapstructure:"status"`
	Endpoint    string    `mapstructure:"endpoint"`
	APIToken    string    `mapstructure:"api_token"`
	TimeoutSec  int       `mapstructure:"timeout_sec"`
	Tools       []MCPTool `mapstructure:"tools"`
}

type MCPTool struct {
	Name         string                 `mapstructure:"name"`
	Description  string                 `mapstructure:"description"`
	Version      string                 `mapstructure:"version"`
	InputSchema  map[string]interface{} `mapstructure:"input_schema"`
	OutputSchema map[string]interface{} `mapstructure:"output_schema"`
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
	Username        string
	Password        string
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

type MinIO struct {
	Endpoint   string
	AccessKey  string
	SecretKey  string
	BucketName string
	UseSSL     bool
	PublicBase string
}

func (s Server) Addr() string {
	return fmt.Sprintf(":%d", s.Port)
}

func Load() (*Config, error) {
	v := viper.New()
	env := strings.ToLower(os.Getenv("GOYAVISION_ENV"))
	if env == "" {
		env = "dev"
	}

	_ = godotenv.Overload("./configs/.env")

	v.SetConfigName(fmt.Sprintf("config.%s", env))
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config for env %s: %w", env, err)
	}

	v.SetEnvPrefix("GOYAVISION")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	cfg := &Config{
		Env:    env,
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
			Username:        v.GetString("mediamtx.username"),
			Password:        v.GetString("mediamtx.password"),
		},
		MinIO: MinIO{
			Endpoint:   v.GetString("minio.endpoint"),
			AccessKey:  v.GetString("minio.access_key"),
			SecretKey:  v.GetString("minio.secret_key"),
			BucketName: v.GetString("minio.bucket_name"),
			UseSSL:     v.GetBool("minio.use_ssl"),
			PublicBase: v.GetString("minio.public_base"),
		},
	}
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
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
	if cfg.AI.Timeout == 0 {
		cfg.AI.Timeout = 10 * time.Second
	}
	if cfg.AI.Retry == 0 {
		cfg.AI.Retry = 2
	}
	if cfg.Record.SegmentSec == 0 {
		cfg.Record.SegmentSec = 300
	}
	if cfg.FFmpeg.MaxRecord == 0 {
		cfg.FFmpeg.MaxRecord = 16
	}
	if cfg.FFmpeg.MaxFrame == 0 {
		cfg.FFmpeg.MaxFrame = 16
	}
	if cfg.Preview.MaxPreview == 0 {
		cfg.Preview.MaxPreview = 10
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	cfg.EncryptKey = v.GetString("encrypt_key")

	_ = v.UnmarshalKey("mcp", &cfg.MCP)
	_ = v.UnmarshalKey("oauth", &cfg.OAuth)
	_ = v.UnmarshalKey("payment", &cfg.Payment)
	return cfg, nil
}

func (c *Config) Validate() error {
	if c.Server.Port == 0 {
		return fmt.Errorf("server.port is required")
	}
	if c.DB.DSN == "" {
		return fmt.Errorf("db.dsn is required")
	}
	if c.JWT.Secret == "" {
		return fmt.Errorf("jwt.secret is required")
	}
	if c.MediaMTX.APIAddress == "" {
		return fmt.Errorf("mediamtx.api_address is required")
	}
	if c.MediaMTX.RTSPAddress == "" || c.MediaMTX.RTMPAddress == "" || c.MediaMTX.HLSAddress == "" || c.MediaMTX.WebRTCAddress == "" || c.MediaMTX.PlaybackAddress == "" {
		return fmt.Errorf("mediamtx addresses are required")
	}
	if c.MediaMTX.RecordPath == "" || c.MediaMTX.RecordFormat == "" || c.MediaMTX.SegmentDuration == "" {
		return fmt.Errorf("mediamtx record settings are required")
	}
	if c.MinIO.Endpoint == "" || c.MinIO.AccessKey == "" || c.MinIO.SecretKey == "" || c.MinIO.BucketName == "" {
		return fmt.Errorf("minio configuration is required")
	}
	return nil
}
