package configs

import (
	"time"

	"github.com/jaydeep/go-n8n/pkg/database"
	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	App        AppConfig        `mapstructure:"app"`
	Server     ServerConfig     `mapstructure:"server"`
	Database   database.Config  `mapstructure:"database"`
	Redis      RedisConfig      `mapstructure:"redis"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Security   SecurityConfig   `mapstructure:"security"`
	CORS       CORSConfig       `mapstructure:"cors"`
	RateLimit  RateLimitConfig  `mapstructure:"rate_limit"`
	Engine     EngineConfig     `mapstructure:"engine"`
	Node       NodeConfig       `mapstructure:"node"`
	Storage    StorageConfig    `mapstructure:"storage"`
	Logging    LoggingConfig    `mapstructure:"logging"`
	Monitoring MonitoringConfig `mapstructure:"monitoring"`
	Webhook    WebhookConfig    `mapstructure:"webhook"`
	Scheduler  SchedulerConfig  `mapstructure:"scheduler"`
	Worker     WorkerConfig     `mapstructure:"worker"`
	Email      EmailConfig      `mapstructure:"email"`
	OAuth      OAuthConfig      `mapstructure:"oauth"`
	Features   FeaturesConfig   `mapstructure:"features"`
	Limits     LimitsConfig     `mapstructure:"limits"`
}

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
	Debug       bool   `mapstructure:"debug"`
}

type ServerConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

type RedisConfig struct {
	Addr         string        `mapstructure:"addr"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	MaxRetries   int           `mapstructure:"max_retries"`
	PoolSize     int           `mapstructure:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns"`
	MaxConnAge   time.Duration `mapstructure:"max_conn_age"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	PoolTimeout  time.Duration `mapstructure:"pool_timeout"`
}

type JWTConfig struct {
	Secret            string        `mapstructure:"secret"`
	AccessTokenExpiry time.Duration `mapstructure:"access_token_expiry"`
	RefreshTokenExpiry time.Duration `mapstructure:"refresh_token_expiry"`
	Issuer            string        `mapstructure:"issuer"`
}

type SecurityConfig struct {
	BCryptCost       int           `mapstructure:"bcrypt_cost"`
	EncryptionKey    string        `mapstructure:"encryption_key"`
	APIKeyLength     int           `mapstructure:"api_key_length"`
	SessionLifetime  time.Duration `mapstructure:"session_lifetime"`
}

type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"allowed_origins"`
	AllowedMethods   []string `mapstructure:"allowed_methods"`
	AllowedHeaders   []string `mapstructure:"allowed_headers"`
	ExposedHeaders   []string `mapstructure:"exposed_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
}

type RateLimitConfig struct {
	Enabled  bool          `mapstructure:"enabled"`
	Requests int           `mapstructure:"requests"`
	Duration time.Duration `mapstructure:"duration"`
	Burst    int           `mapstructure:"burst"`
}

type EngineConfig struct {
	MaxParallelExecutions int           `mapstructure:"max_parallel_executions"`
	MaxExecutionTime      time.Duration `mapstructure:"max_execution_time"`
	WorkerCount          int           `mapstructure:"worker_count"`
	QueueSize            int           `mapstructure:"queue_size"`
	MaxRetries           int           `mapstructure:"max_retries"`
	RetryBackoff         time.Duration `mapstructure:"retry_backoff"`
	CheckpointInterval   time.Duration `mapstructure:"checkpoint_interval"`
}

type NodeConfig struct {
	MaxExecutionTime      time.Duration `mapstructure:"max_execution_time"`
	EnableDynamicLoading  bool          `mapstructure:"enable_dynamic_loading"`
	SandboxExecution      bool          `mapstructure:"sandbox_execution"`
	MaxDataSize          int64         `mapstructure:"max_data_size"`
	Timeout              time.Duration `mapstructure:"timeout"`
}

type StorageConfig struct {
	Type  string              `mapstructure:"type"`
	Local LocalStorageConfig  `mapstructure:"local"`
	S3    S3StorageConfig     `mapstructure:"s3"`
}

type LocalStorageConfig struct {
	Path string `mapstructure:"path"`
}

type S3StorageConfig struct {
	Bucket    string `mapstructure:"bucket"`
	Region    string `mapstructure:"region"`
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
}

type LoggingConfig struct {
	Level  string         `mapstructure:"level"`
	Format string         `mapstructure:"format"`
	Output string         `mapstructure:"output"`
	File   FileLogConfig  `mapstructure:"file"`
}

type FileLogConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	Path       string `mapstructure:"path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
	Compress   bool   `mapstructure:"compress"`
}

type MonitoringConfig struct {
	Metrics MetricsConfig `mapstructure:"metrics"`
	Health  HealthConfig  `mapstructure:"health"`
	Tracing TracingConfig `mapstructure:"tracing"`
}

type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Port    int    `mapstructure:"port"`
	Path    string `mapstructure:"path"`
}

type HealthConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Path     string `mapstructure:"path"`
	Detailed bool   `mapstructure:"detailed"`
}

type TracingConfig struct {
	Enabled        bool   `mapstructure:"enabled"`
	ServiceName    string `mapstructure:"service_name"`
	JaegerEndpoint string `mapstructure:"jaeger_endpoint"`
}

type WebhookConfig struct {
	BaseURL         string        `mapstructure:"base_url"`
	Timeout         time.Duration `mapstructure:"timeout"`
	MaxPayloadSize  int64         `mapstructure:"max_payload_size"`
	RetryAttempts   int           `mapstructure:"retry_attempts"`
	RetryDelay      time.Duration `mapstructure:"retry_delay"`
}

type SchedulerConfig struct {
	Enabled           bool          `mapstructure:"enabled"`
	CheckInterval     time.Duration `mapstructure:"check_interval"`
	Location          string        `mapstructure:"location"`
	MaxConcurrentJobs int           `mapstructure:"max_concurrent_jobs"`
}

type WorkerConfig struct {
	Concurrency     int           `mapstructure:"concurrency"`
	QueueName       string        `mapstructure:"queue_name"`
	RetryMax        int           `mapstructure:"retry_max"`
	RetryDelay      time.Duration `mapstructure:"retry_delay"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

type EmailConfig struct {
	Enabled bool       `mapstructure:"enabled"`
	SMTP    SMTPConfig `mapstructure:"smtp"`
}

type SMTPConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
	UseTLS   bool   `mapstructure:"use_tls"`
}

type OAuthConfig struct {
	Google OAuthProviderConfig `mapstructure:"google"`
	GitHub OAuthProviderConfig `mapstructure:"github"`
}

type OAuthProviderConfig struct {
	Enabled      bool     `mapstructure:"enabled"`
	ClientID     string   `mapstructure:"client_id"`
	ClientSecret string   `mapstructure:"client_secret"`
	RedirectURL  string   `mapstructure:"redirect_url"`
	Scopes       []string `mapstructure:"scopes"`
}

type FeaturesConfig struct {
	Teams         bool `mapstructure:"teams"`
	Marketplace   bool `mapstructure:"marketplace"`
	CustomNodes   bool `mapstructure:"custom_nodes"`
	WebhookTunnel bool `mapstructure:"webhook_tunnel"`
	APIAccess     bool `mapstructure:"api_access"`
	OAuthLogin    bool `mapstructure:"oauth_login"`
	TwoFactorAuth bool `mapstructure:"two_factor_auth"`
}

type LimitsConfig struct {
	MaxWorkflowsPerUser      int           `mapstructure:"max_workflows_per_user"`
	MaxNodesPerWorkflow      int           `mapstructure:"max_nodes_per_workflow"`
	MaxExecutionTime         time.Duration `mapstructure:"max_execution_time"`
	MaxFileSize              int64         `mapstructure:"max_file_size"`
	MaxAPIRequestsPerMinute  int           `mapstructure:"max_api_requests_per_minute"`
}

// Load loads configuration from file and environment
func Load() (*Config, error) {
	viper.SetConfigFile("configs/config.yaml")
	viper.SetConfigType("yaml")
	
	// Read from environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("N8N")
	
	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	
	// Override with environment variables
	loadEnvOverrides(&config)
	
	return &config, nil
}

// loadEnvOverrides loads environment variable overrides
func loadEnvOverrides(cfg *Config) {
	// Override critical settings from environment
	if viper.IsSet("DB_HOST") {
		cfg.Database.Host = viper.GetString("DB_HOST")
	}
	if viper.IsSet("DB_PASSWORD") {
		cfg.Database.Password = viper.GetString("DB_PASSWORD")
	}
	if viper.IsSet("REDIS_URL") {
		cfg.Redis.Addr = viper.GetString("REDIS_URL")
	}
	if viper.IsSet("JWT_SECRET") {
		cfg.JWT.Secret = viper.GetString("JWT_SECRET")
	}
	if viper.IsSet("ENCRYPTION_KEY") {
		cfg.Security.EncryptionKey = viper.GetString("ENCRYPTION_KEY")
	}
}
