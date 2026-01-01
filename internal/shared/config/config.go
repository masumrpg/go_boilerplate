package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	Redis      RedisConfig
	JWT        JWTConfig
	OAuth      OAuthConfig
	Email      EmailConfig
	Security   SecurityConfig
	Logger     LoggerConfig
	SuperAdmin SuperAdminConfig
}

// SecurityConfig holds security configuration
type SecurityConfig struct {
	EmailVerificationEnabled bool `mapstructure:"EMAIL_VERIFICATION_ENABLED"`
	TwoFactorEnabled         bool `mapstructure:"TWO_FACTOR_ENABLED"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string `mapstructure:"SERVER_PORT"`
	Host string `mapstructure:"SERVER_HOST"`
	Mode string `mapstructure:"SERVER_MODE"` // development, production, test
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	DBName   string `mapstructure:"DB_NAME"`
	SSLMode  string `mapstructure:"DB_SSLMODE"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     string `mapstructure:"REDIS_PORT"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret          string `mapstructure:"JWT_SECRET"`
	AccessExpiry    time.Duration
	RefreshExpiry   time.Duration
	Issuer          string
}

// OAuthConfig holds OAuth configuration
type OAuthConfig struct {
	Google GoogleOAuthConfig
	GitHub GitHubOAuthConfig
}

// GoogleOAuthConfig holds Google OAuth configuration
type GoogleOAuthConfig struct {
	ClientID         string `mapstructure:"OAUTH_GOOGLE_CLIENT_ID"`
	ClientSecret     string `mapstructure:"OAUTH_GOOGLE_CLIENT_SECRET"`
	RedirectURL      string `mapstructure:"OAUTH_GOOGLE_REDIRECT_URL"`
	Enabled          bool   `mapstructure:"OAUTH_GOOGLE_ENABLED"`
	SendWelcomeEmail bool   `mapstructure:"OAUTH_GOOGLE_SEND_WELCOME_EMAIL"`
}

// GitHubOAuthConfig holds GitHub OAuth configuration
type GitHubOAuthConfig struct {
	ClientID         string `mapstructure:"OAUTH_GITHUB_CLIENT_ID"`
	ClientSecret     string `mapstructure:"OAUTH_GITHUB_CLIENT_SECRET"`
	RedirectURL      string `mapstructure:"OAUTH_GITHUB_REDIRECT_URL"`
	Enabled          bool   `mapstructure:"OAUTH_GITHUB_ENABLED"`
	SendWelcomeEmail bool   `mapstructure:"OAUTH_GITHUB_SEND_WELCOME_EMAIL"`
}

// EmailConfig holds email configuration
type EmailConfig struct {
	SMTPHost     string `mapstructure:"SMTP_HOST"`
	SMTPPort     int    `mapstructure:"SMTP_PORT"`
	SMTPUser     string `mapstructure:"SMTP_USER"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD"`
	SMTPFrom     string `mapstructure:"SMTP_FROM"`
	Enabled      bool   `mapstructure:"EMAIL_ENABLED"`
}

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	Level  string `mapstructure:"LOG_LEVEL"` // debug, info, warn, error
	Format string `mapstructure:"LOG_FORMAT"` // json, text
}

// SuperAdminConfig holds default SuperAdmin account configuration
type SuperAdminConfig struct {
	Email    string `mapstructure:"SUPERADMIN_EMAIL"`
	Password string `mapstructure:"SUPERADMIN_PASSWORD"`
	Name     string `mapstructure:"SUPERADMIN_NAME"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		fmt.Println("⚠️  No .env file found, using system environment variables or defaults")
	} else {
		fmt.Println("✅ .env file loaded successfully")
	}

	// Set defaults
	setDefaults()

	// Debug: Show what env vars we're reading
	fmt.Println("📖 Reading environment variables...")

	// Create config from environment variables
	cfg := Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "3000"),
			Host: getEnv("SERVER_HOST", "localhost"),
			Mode: getEnv("SERVER_MODE", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "go_boilerplate"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       parseInt(getEnv("REDIS_DB", "0")),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", ""),
		},
		OAuth: OAuthConfig{
			Google: GoogleOAuthConfig{
				ClientID:         getEnv("OAUTH_GOOGLE_CLIENT_ID", ""),
				ClientSecret:     getEnv("OAUTH_GOOGLE_CLIENT_SECRET", ""),
				RedirectURL:      getEnv("OAUTH_GOOGLE_REDIRECT_URL", ""),
				Enabled:          getBoolEnv("OAUTH_GOOGLE_ENABLED", false),
				SendWelcomeEmail: getBoolEnv("OAUTH_GOOGLE_SEND_WELCOME_EMAIL", false),
			},
			GitHub: GitHubOAuthConfig{
				ClientID:         getEnv("OAUTH_GITHUB_CLIENT_ID", ""),
				ClientSecret:     getEnv("OAUTH_GITHUB_CLIENT_SECRET", ""),
				RedirectURL:      getEnv("OAUTH_GITHUB_REDIRECT_URL", ""),
				Enabled:          getBoolEnv("OAUTH_GITHUB_ENABLED", false),
				SendWelcomeEmail: getBoolEnv("OAUTH_GITHUB_SEND_WELCOME_EMAIL", false),
			},
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:     parseInt(getEnv("SMTP_PORT", "587")),
			SMTPUser:     getEnv("SMTP_USER", ""),
			SMTPPassword: getEnv("SMTP_PASSWORD", ""),
			SMTPFrom:     getEnv("SMTP_FROM", ""),
			Enabled:      getBoolEnv("EMAIL_ENABLED", false),
		},
		Security: SecurityConfig{
			EmailVerificationEnabled: getBoolEnv("EMAIL_VERIFICATION_ENABLED", false),
			TwoFactorEnabled:         getBoolEnv("TWO_FACTOR_ENABLED", false),
		},
		Logger: LoggerConfig{
			Level:  getEnv("LOG_LEVEL", "debug"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		SuperAdmin: SuperAdminConfig{
			Name:     getEnv("SUPERADMIN_NAME", "Super Admin"),
			Email:    getEnv("SUPERADMIN_EMAIL", "superadmin@boilerplate.com"),
			Password: getEnv("SUPERADMIN_PASSWORD", "SuperAdmin123!"),
		},
	}

	// Parse JWT expiry durations
	var err error
	cfg.JWT.AccessExpiry, err = time.ParseDuration(getEnv("JWT_ACCESS_EXPIRY", "1h"))
	if err != nil {
		cfg.JWT.AccessExpiry = 1 * time.Hour
	}

	cfg.JWT.RefreshExpiry, err = time.ParseDuration(getEnv("JWT_REFRESH_EXPIRY", "24h"))
	if err != nil {
		cfg.JWT.RefreshExpiry = 24 * time.Hour
	}

	cfg.JWT.Issuer = "go_boilerplate"

	// Debug: Print loaded config
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📋 Configuration Loaded:")
	fmt.Printf("   Server Port: %s\n", cfg.Server.Port)
	fmt.Printf("   Server Mode: %s\n", cfg.Server.Mode)
	fmt.Printf("   Database: %s@%s:%s/%s\n", cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	fmt.Printf("   Redis: %s:%s (DB: %d)\n", cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.DB)
	fmt.Printf("   JWT Secret: %s\n", maskSecret(cfg.JWT.Secret))
	fmt.Printf("   Log Level: %s\n", cfg.Logger.Level)
	fmt.Printf("   Security: EmailVerify=%v, 2FA=%v\n", cfg.Security.EmailVerificationEnabled, cfg.Security.TwoFactorEnabled)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Validate required fields
	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

// maskSecret masks JWT secret for display
func maskSecret(secret string) string {
	if secret == "" {
		return "[NOT SET]"
	}
	if len(secret) > 8 {
		return secret[:4] + "****" + secret[len(secret)-4:]
	}
	return "****"
}

// getEnv gets an environment variable or returns the default value
func getEnv(key, defaultValue string) string {
	// Try os.Getenv first (from godotenv)
	if value := os.Getenv(key); value != "" {
		fmt.Printf("   ✅ %s = %s (from .env)\n", key, value)
		return value
	}

	// Fallback to viper
	if value := viper.GetString(key); value != "" {
		fmt.Printf("   ✅ %s = %s (from system)\n", key, value)
		return value
	}

	// Use default
	fmt.Printf("   ⚠️  %s not set, using default: %s\n", key, defaultValue)
	return defaultValue
}

// parseInt parses a string to int
func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

// getBoolEnv parses a string to bool
func getBoolEnv(key string, defaultValue bool) bool {
	// Try os.Getenv first (from godotenv)
	if value := os.Getenv(key); value != "" {
	 parsed := parseBool(value)
	 fmt.Printf("   ✅ %s = %v (from .env)\n", key, parsed)
	 return parsed
	}

	// Fallback to viper
	if value := viper.GetString(key); value != "" {
	 parsed := parseBool(value)
	 fmt.Printf("   ✅ %s = %v (from system)\n", key, parsed)
	 return parsed
	}

	// Use default
	fmt.Printf("   ⚠️  %s not set, using default: %v\n", key, defaultValue)
	return defaultValue
}

// parseBool parses a string to bool (accepts: true, false, 1, 0, yes, no)
func parseBool(s string) bool {
	switch s {
	case "true", "1", "yes", "TRUE", "YES", "True":
		return true
	case "false", "0", "no", "FALSE", "NO", "False":
		return false
	default:
		return false
	}
}

// bindEnvs binds environment variables to config keys
func bindEnvs() {
	viper.BindEnv("SERVER_PORT")
	viper.BindEnv("SERVER_HOST")
	viper.BindEnv("SERVER_MODE")

	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_PORT")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_NAME")
	viper.BindEnv("DB_SSLMODE")

	viper.BindEnv("JWT_SECRET")
	viper.BindEnv("JWT_ACCESS_EXPIRY")
	viper.BindEnv("JWT_REFRESH_EXPIRY")

	viper.BindEnv("OAUTH_GOOGLE_CLIENT_ID")
	viper.BindEnv("OAUTH_GOOGLE_CLIENT_SECRET")
	viper.BindEnv("OAUTH_GOOGLE_REDIRECT_URL")

	viper.BindEnv("OAUTH_GITHUB_CLIENT_ID")
	viper.BindEnv("OAUTH_GITHUB_CLIENT_SECRET")
	viper.BindEnv("OAUTH_GITHUB_REDIRECT_URL")

	viper.BindEnv("SMTP_HOST")
	viper.BindEnv("SMTP_PORT")
	viper.BindEnv("SMTP_USER")
	viper.BindEnv("SMTP_PASSWORD")
	viper.BindEnv("SMTP_FROM")

	viper.BindEnv("EMAIL_VERIFICATION_ENABLED")
	viper.BindEnv("TWO_FACTOR_ENABLED")

	viper.BindEnv("LOG_LEVEL")
	viper.BindEnv("LOG_FORMAT")
}

// setDefaults sets default configuration values
func setDefaults() {
	// Server defaults
	viper.SetDefault("SERVER_PORT", "3000")
	viper.SetDefault("SERVER_HOST", "localhost")
	viper.SetDefault("SERVER_MODE", "development")

	// Database defaults
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "go_boilerplate")
	viper.SetDefault("DB_SSLMODE", "disable")

	// JWT defaults
	viper.SetDefault("JWT_SECRET", "change-this-secret-in-production")
	viper.SetDefault("JWT_ACCESS_EXPIRY", "1h")
	viper.SetDefault("JWT_REFRESH_EXPIRY", "24h")

	// Email defaults
	viper.SetDefault("SMTP_PORT", "587")

	// Security defaults
	viper.SetDefault("EMAIL_VERIFICATION_ENABLED", false)
	viper.SetDefault("TWO_FACTOR_ENABLED", false)

	// Logger defaults
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("LOG_FORMAT", "json")
}

// validateConfig validates required configuration fields
func validateConfig(cfg *Config) error {
	if cfg.Server.Port == "" {
		return fmt.Errorf("SERVER_PORT is required")
	}
	if cfg.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if cfg.Database.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	// Only require JWT_SECRET in production
	if cfg.Server.IsProduction() && (cfg.JWT.Secret == "" || cfg.JWT.Secret == "change-this-secret-in-production") {
		return fmt.Errorf("JWT_SECRET must be set to a secure value in production")
	}
	// In development, use a default secret if not set
	if cfg.JWT.Secret == "" && cfg.Server.IsDevelopment() {
		cfg.JWT.Secret = "development-secret-key-change-in-production"
		fmt.Println("WARNING: Using default JWT secret for development mode!")
	}
	return nil
}

// GetDSN returns the PostgreSQL Data Source Name
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// IsDevelopment returns true if server mode is development
func (c *ServerConfig) IsDevelopment() bool {
	return c.Mode == "development"
}

// IsProduction returns true if server mode is production
func (c *ServerConfig) IsProduction() bool {
	return c.Mode == "production"
}
