package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

var (
	instance *Config
	once     sync.Once
	mu       sync.RWMutex
)

// Config holds all application configuration
type Config struct {
	MaintenanceEligibleFlatStatuses []string
	Backup                          BackupConfig
	ImageKit                        ImageKitConfig
	// Application
	Environment  string
	AppName      string
	Version      string
	Host         string
	Port         string
	PublicAppURL string

	// Logging
	LogLevel      string
	LogDir        string
	LogMaxSize    int // MB
	LogMaxBackups int
	LogMaxAge     int // days
	LogCompress   bool
	EnableConsole bool

	// Server timeouts
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	// CORS
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool

	// Rate limiting
	RateLimitEnabled bool
	RateLimitRPS     int
	RateLimitBurst   int

	// Database
	DBHost         string
	DBPort         int
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	DBMaxOpenConns int
	DBMaxIdleConns int
	DBMaxLifetime  time.Duration
	AutoMigrate    bool

	// Redis (optional)
	RedisEnabled  bool
	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisDB       int

	// JWT
	JWTSecret             string
	JWTAccessTokenExpiry  time.Duration
	JWTRefreshTokenExpiry time.Duration
	JWTIssuer             string

	// OTP
	OTPSecret string
	OTPExpiry time.Duration

	// Auth groups security settings used by auth services and handlers.
	Auth AuthConfig

	// Email
	ResendAPIKey  string
	EmailFrom     string
	EmailFromName string

	// Push notifications (FCM)
	FCMEnabled         bool
	FCMWebEnabled      bool
	FCMCredentialsPath string
	APNSEnabled bool
	APNSKeyPath string
	APNSKeyID string
	APNSTeamID string
	APNSTopic string
	APNSProduction bool

	// Flat member invites
	FlatMemberInviteTTL time.Duration

	// Background jobs
	JobTimezone                  string
	JobWebhookSecret             string
	VisitorEntryRetentionDays    int
	VisitorInviteRetentionDays   int
	FlatInviteRetentionDays      int
	NotificationRetentionDays    int
	CleanupBatchSize             int
	VisitorReportBatchSize       int
	VisitorReportProcessingLease time.Duration
	VisitorReportAllowResend     bool

	// Security
	EnableHTTPS    bool
	TLSCertFile    string
	TLSKeyFile     string
	TrustedProxies []string
	EnableCSRF     bool
	CSRFSecret     string

	// Monitoring
	EnableMetrics   bool
	MetricsPort     string
	EnableHealthz   bool
	EnableProfiling bool
}

type AuthConfig struct {
	JWTSecret        string
	JWTIssuer        string
	AccessExpiry     time.Duration
	RefreshExpiry    time.Duration
	OnboardingExpiry time.Duration
	IsProduction     bool
	OTPSecret        string
	OTPExpiry        time.Duration
	CookieDomain     string
}

// LoadConfig loads configuration from environment files and variables
func LoadConfig() (*Config, error) {
	var loadErr error
	once.Do(func() {
		instance, loadErr = loadConfigInternal()
	})
	return instance, loadErr
}

// GetConfig returns the singleton config instance (must call LoadConfig first)
func GetConfig() *Config {
	mu.RLock()
	defer mu.RUnlock()
	if instance == nil {
		panic("config not initialized, call LoadConfig() first")
	}
	return instance
}

// MustLoadConfig loads config or panics
func MustLoadConfig() *Config {
	cfg, err := LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}
	return cfg
}

func loadConfigInternal() (*Config, error) {
	// Determine environment
	env := getEnv("GO_ENV", "development")
	env = strings.ToLower(env)

	// Load appropriate .env file
	if err := loadEnvFile(env); err != nil {
		// Non-fatal, continue with system env vars
		fmt.Printf("Warning: %v\n", err)
	}

	jwtSecret := getEnv("JWT_SECRET", "")
	jwtIssuer := getEnv("JWT_ISSUER", "go-server")
	accessExpiry := time.Duration(getEnvAsInt("JWT_ACCESS_TOKEN_EXPIRY", 15*60)) * time.Second
	refreshExpiry := time.Duration(getEnvAsInt("JWT_REFRESH_TOKEN_EXPIRY", 24*60*60*7)) * time.Second
	onboardingExpiry := time.Duration(getEnvAsInt("JWT_ONBOARDING_TOKEN_EXPIRY", 30*60)) * time.Second
	otpSecret := getEnv("OTP_SECRET", "Something_Secret_ChangeMe")
	otpExpiry := time.Duration(getEnvAsInt("OTP_EXPIRY_SECONDS", 10*60)) * time.Second

	flatMemberInviteTTL, err := getEnvAsDuration("FLAT_MEMBER_INVITE_TTL", 168*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("invalid FLAT_MEMBER_INVITE_TTL: %w", err)
	}

	appName := getEnv("APP_NAME", "go-server")

	// Build config
	config := &Config{
		ImageKit: ImageKitConfig{
			Enabled:      getEnvAsBool("IMAGEKIT_ENABLED", false),
			PrivateKey:   getEnv("IMAGEKIT_PRIVATE_KEY", ""),
			URLEndpoint:  getEnv("IMAGEKIT_URL_ENDPOINT", ""),
			FolderPrefix: getEnv("IMAGEKIT_FOLDER_PREFIX", "dev/apna-gate"),
		},
		Backup: loadBackupConfig(),
		// Application
		Environment:  env,
		AppName:      appName,
		Version:      getEnv("APP_VERSION", "1.0.0"),
		Host:         getEnvFirst([]string{"SERVER_HOST", "HOST"}, "127.0.0.1"),
		Port:         getEnvFirst([]string{"SERVER_PORT", "PORT"}, "8080"),
		PublicAppURL: strings.TrimRight(getEnv("PUBLIC_APP_URL", "https://apnagate.app"), "/"),

		// Logging
		LogLevel:      getEnv("LOG_LEVEL", getDefaultLogLevel(env)),
		LogDir:        getEnv("LOG_DIR", "logs"),
		LogMaxSize:    getEnvAsInt("LOG_MAX_SIZE", 100),
		LogMaxBackups: getEnvAsInt("LOG_MAX_BACKUPS", 5),
		LogMaxAge:     getEnvAsInt("LOG_MAX_AGE", 30),
		LogCompress:   getEnvAsBool("LOG_COMPRESS", true),
		EnableConsole: getEnvAsBool("LOG_CONSOLE", env != "production"),

		// Server timeouts (convert seconds to duration)
		ReadTimeout:  time.Duration(getEnvAsInt("READ_TIMEOUT", 60)) * time.Second,
		WriteTimeout: time.Duration(getEnvAsInt("WRITE_TIMEOUT", 90)) * time.Second,
		IdleTimeout:  time.Duration(getEnvAsInt("IDLE_TIMEOUT", 120)) * time.Second,

		// CORS
		AllowOrigins:     getEnvAsSlice("ALLOW_ORIGINS", []string{"*"}),
		AllowMethods:     getEnvAsSlice("ALLOW_METHODS", []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}),
		AllowHeaders:     getEnvAsSlice("ALLOW_HEADERS", []string{"Origin", "Content-Type", "Authorization", "Idempotency-Key"}),
		AllowCredentials: getEnvAsBool("ALLOW_CREDENTIALS", true),

		// Rate limiting
		RateLimitEnabled: getEnvAsBool("RATE_LIMIT_ENABLED", env == "production"),
		RateLimitRPS:     getEnvAsInt("RATE_LIMIT_RPS", 100),
		RateLimitBurst:   getEnvAsInt("RATE_LIMIT_BURST", 200),

		// Database
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnvAsInt("DB_PORT", 5432),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBName:         getEnv("DB_NAME", "app_db"),
		DBSSLMode:      getEnv("DB_SSL_MODE", getDefaultSSLMode(env)),
		DBMaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
		DBMaxLifetime:  time.Duration(getEnvAsInt("DB_MAX_LIFETIME", 300)) * time.Second,
		AutoMigrate:    getEnvAsBool("AUTO_MIGRATE", env != "production"),

		// Redis
		// RedisEnabled:  getEnvAsBool("REDIS_ENABLED", false),
		// RedisHost:     getEnv("REDIS_HOST", "localhost"),
		// RedisPort:     getEnvAsInt("REDIS_PORT", 6379),
		// RedisPassword: getEnv("REDIS_PASSWORD", ""),
		// RedisDB:       getEnvAsInt("REDIS_DB", 0),

		// JWT
		JWTSecret: jwtSecret,
		// JWTSecret:        getEnv("JWT_SECRET", ""),
		// JWTExpiry:        time.Duration(getEnvAsInt("JWT_EXPIRY", 3600)) * time.Hour,
		// JWTRefreshExpiry: time.Duration(getEnvAsInt("JWT_REFRESH_EXPIRY", 86400)) * time.Hour,
		JWTAccessTokenExpiry:  accessExpiry,
		JWTRefreshTokenExpiry: refreshExpiry,
		JWTIssuer:             jwtIssuer,

		// OTP
		OTPSecret: otpSecret,
		OTPExpiry: otpExpiry,

		// Auth
		Auth: AuthConfig{
			JWTSecret:        jwtSecret,
			JWTIssuer:        jwtIssuer,
			AccessExpiry:     accessExpiry,
			RefreshExpiry:    refreshExpiry,
			OnboardingExpiry: onboardingExpiry,
			IsProduction:     env == "production" || env == "prod",
			OTPSecret:        otpSecret,
			OTPExpiry:        otpExpiry,
			CookieDomain:     getEnv("AUTH_COOKIE_DOMAIN", ""),
		},

		// Email
		ResendAPIKey:  getEnv("RESEND_API_KEY", ""),
		EmailFrom:     getEnv("EMAIL_FROM", ""),
		EmailFromName: getEnv("EMAIL_FROM_NAME", appName),

		// Push notifications
		FCMEnabled:         getEnvAsBool("FCM_ENABLED", false),
		FCMWebEnabled:      getEnvAsBool("FCM_WEB_ENABLED", false),
		FCMCredentialsPath: getEnv("FCM_CREDENTIALS_PATH", ""),
		APNSEnabled: getEnvAsBool("APNS_ENABLED", false),
		APNSKeyPath: getEnv("APNS_KEY_PATH", ""),
		APNSKeyID: getEnv("APNS_KEY_ID", ""),
		APNSTeamID: getEnv("APNS_TEAM_ID", ""),
		APNSTopic: getEnv("APNS_TOPIC", ""),
		APNSProduction: getEnvAsBool("APNS_PRODUCTION", false),

		FlatMemberInviteTTL: flatMemberInviteTTL,

		// Background jobs
		MaintenanceEligibleFlatStatuses: maintenanceStatusesFromEnv(),
		JobTimezone:                     getEnv("JOB_TIMEZONE", "Asia/Kolkata"),
		JobWebhookSecret:                strings.TrimSpace(getEnv("JOB_WEBHOOK_SECRET", "")),
		VisitorEntryRetentionDays:       getEnvAsInt("VISITOR_ENTRY_RETENTION_DAYS", 90),
		VisitorInviteRetentionDays:      getEnvAsInt("VISITOR_INVITE_RETENTION_DAYS", 90),
		FlatInviteRetentionDays:         getEnvAsInt("FLAT_MEMBER_INVITE_RETENTION_DAYS", 90),
		NotificationRetentionDays:       getEnvAsInt("NOTIFICATION_RETENTION_DAYS", 30),
		CleanupBatchSize:                getEnvAsInt("CLEANUP_BATCH_SIZE", 500),
		VisitorReportBatchSize:          getEnvAsInt("VISITOR_REPORT_BATCH_SIZE", 50),
		VisitorReportProcessingLease:    15 * time.Minute,
		VisitorReportAllowResend:        getEnvAsBool("VISITOR_REPORT_ALLOW_RESEND", false),

		// Security
		// EnableHTTPS:    getEnvAsBool("ENABLE_HTTPS", env == "production"),
		// TLSCertFile:    getEnv("TLS_CERT_FILE", ""),
		// TLSKeyFile:     getEnv("TLS_KEY_FILE", ""),
		// TrustedProxies: getEnvAsSlice("TRUSTED_PROXIES", []string{}),
		// EnableCSRF:     getEnvAsBool("ENABLE_CSRF", env == "production"),
		// CSRFSecret:     getEnv("CSRF_SECRET", ""),

		// Monitoring
		EnableMetrics:   getEnvAsBool("ENABLE_METRICS", true),
		MetricsPort:     getEnv("METRICS_PORT", "9090"),
		EnableHealthz:   getEnvAsBool("ENABLE_HEALTHZ", true),
		EnableProfiling: getEnvAsBool("ENABLE_PROFILING", env != "production"),
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	fmt.Printf("Configuration loaded successfully [env=%s]\n", config.Environment)

	return config, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if _, err := c.MaintenanceStatuses(); err != nil {
		return err
	}
	if err := c.Backup.Validate(); err != nil {
		return err
	}
	if err := c.validateImages(); err != nil {
		return err
	}
	if c.VisitorReportAllowResend && !c.IsDevelopment() {
		return fmt.Errorf("VISITOR_REPORT_ALLOW_RESEND may only be enabled in development")
	}
	// Production-specific validations
	if c.Environment == "production" {
		if c.DBPassword == "" {
			return fmt.Errorf("DB_PASSWORD is required in production")
		}
		if c.JWTSecret == "" {
			return fmt.Errorf("JWT_SECRET is required in production")
		}
		if c.OTPSecret == "" {
			return fmt.Errorf("OTP_SECRET is required in production")
		}
		if c.EnableHTTPS && (c.TLSCertFile == "" || c.TLSKeyFile == "") {
			return fmt.Errorf("TLS_CERT_FILE and TLS_KEY_FILE required when HTTPS is enabled")
		}
		if len(c.AllowOrigins) == 0 {
			return fmt.Errorf("ALLOW_ORIGINS is required in production")
		}
		if c.OTPSecret == "" {
			return fmt.Errorf("OTP_SECRET is required in production")
		}
		for _, origin := range c.AllowOrigins {
			if origin == "*" {
				return fmt.Errorf("wildcard CORS origins not allowed in production")
			}
		}
		if !c.FCMEnabled {
			return fmt.Errorf("FCM_ENABLED must be true in production")
		}
		if strings.TrimSpace(c.FCMCredentialsPath) == "" {
			return fmt.Errorf("FCM_CREDENTIALS_PATH is required in production")
		}
	}

	// General validations
	if c.Port == "" {
		return fmt.Errorf("PORT is required")
	}
	if c.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if c.ReadTimeout < 0 || c.WriteTimeout < 0 || c.IdleTimeout < 0 {
		return fmt.Errorf("timeouts cannot be negative")
	}
	if c.OTPExpiry <= 0 {
		return fmt.Errorf("OTP_EXPIRY_SECONDS must be positive")
	}
	if c.Auth.OnboardingExpiry <= 0 {
		return fmt.Errorf("JWT_ONBOARDING_TOKEN_EXPIRY must be positive")
	}
	if _, err := time.LoadLocation(c.JobTimezone); err != nil {
		return fmt.Errorf("JOB_TIMEZONE is invalid: %w", err)
	}
	if len(strings.TrimSpace(c.JobWebhookSecret)) < 32 {
		return fmt.Errorf("JOB_WEBHOOK_SECRET must contain at least 32 characters")
	}
	if c.VisitorEntryRetentionDays <= 0 || c.VisitorInviteRetentionDays <= 0 ||
		c.FlatInviteRetentionDays <= 0 || c.NotificationRetentionDays <= 0 {
		return fmt.Errorf("job retention days must be positive")
	}
	if c.CleanupBatchSize <= 0 || c.VisitorReportBatchSize <= 0 {
		return fmt.Errorf("job batch sizes must be positive")
	}
	if c.CleanupBatchSize > 100_000 || c.VisitorReportBatchSize > 100_000 {
		return fmt.Errorf("job batch sizes must not exceed 100000")
	}
	if c.VisitorReportProcessingLease <= 0 {
		return fmt.Errorf("visitor report processing lease must be positive")
	}

	return nil
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development" || c.Environment == "dev"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production" || c.Environment == "prod"
}

// IsTesting returns true if running in test mode
func (c *Config) IsTesting() bool {
	return c.Environment == "test" || c.Environment == "testing"
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

// GetRedisAddr returns Redis address
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.RedisHost, c.RedisPort)
}

// GetServerAddr returns the server address
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// ==================== Helper Functions ====================

func loadEnvFile(env string) error {
	if envFile := strings.TrimSpace(os.Getenv("ENV_FILE")); envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			return fmt.Errorf("failed to load ENV_FILE %q: %w", envFile, err)
		}
		fmt.Printf("Loaded %s\n", envFile)
		return nil
	}

	// Try environment-specific file first
	envFile := fmt.Sprintf(".env.%s", env)
	if err := godotenv.Load(envFile); err == nil {
		fmt.Printf("Loaded %s\n", envFile)
		return nil
	}

	// Try default .env file
	if err := godotenv.Load(); err == nil {
		fmt.Printf("Loaded .env\n")
		return nil
	}

	return fmt.Errorf("no .env file found, using system environment variables")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvFirst(keys []string, defaultValue string) string {
	for _, key := range keys {
		if value := os.Getenv(key); value != "" {
			return value
		}
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	fmt.Printf("Warning: Invalid integer value for %s, using default %d\n", key, defaultValue)
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		fmt.Printf("Warning: Invalid boolean value for %s, using default %v\n", key, defaultValue)
		return defaultValue
	}
	return value
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	// Split by comma and trim spaces
	values := strings.Split(valueStr, ",")
	for i, v := range values {
		values[i] = strings.TrimSpace(v)
	}
	return values
}

func getEnvAsDuration(key string, defaultValue time.Duration) (time.Duration, error) {
	valueStr := strings.TrimSpace(os.Getenv(key))
	if valueStr == "" {
		return defaultValue, nil
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", key, err)
	}
	if value <= 0 {
		return 0, fmt.Errorf("%s must be positive", key)
	}
	return value, nil
}

func getDefaultLogLevel(env string) string {
	if env == "production" {
		return "info"
	}
	return "debug"
}

func getDefaultSSLMode(env string) string {
	if env == "production" {
		return "require"
	}
	return "disable"
}
