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
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	OAuth    OAuthConfig
	Email    EmailConfig
	Logger   LoggerConfig
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
	ClientID     string `mapstructure:"OAUTH_GOOGLE_CLIENT_ID"`
	ClientSecret string `mapstructure:"OAUTH_GOOGLE_CLIENT_SECRET"`
	RedirectURL  string `mapstructure:"OAUTH_GOOGLE_REDIRECT_URL"`
}

// GitHubOAuthConfig holds GitHub OAuth configuration
type GitHubOAuthConfig struct {
	ClientID     string `mapstructure:"OAUTH_GITHUB_CLIENT_ID"`
	ClientSecret string `mapstructure:"OAUTH_GITHUB_CLIENT_SECRET"`
	RedirectURL  string `mapstructure:"OAUTH_GITHUB_REDIRECT_URL"`
}

// EmailConfig holds email configuration
type EmailConfig struct {
	SMTPHost     string `mapstructure:"SMTP_HOST"`
	SMTPPort     int    `mapstructure:"SMTP_PORT"`
	SMTPUser     string `mapstructure:"SMTP_USER"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD"`
	SMTPFrom     string `mapstructure:"SMTP_FROM"`
}

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	Level  string `mapstructure:"LOG_LEVEL"` // debug, info, warn, error
	Format string `mapstructure:"LOG_FORMAT"` // json, text
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
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", ""),
		},
		OAuth: OAuthConfig{
			Google: GoogleOAuthConfig{
				ClientID:     getEnv("OAUTH_GOOGLE_CLIENT_ID", ""),
				ClientSecret: getEnv("OAUTH_GOOGLE_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("OAUTH_GOOGLE_REDIRECT_URL", ""),
			},
			GitHub: GitHubOAuthConfig{
				ClientID:     getEnv("OAUTH_GITHUB_CLIENT_ID", ""),
				ClientSecret: getEnv("OAUTH_GITHUB_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("OAUTH_GITHUB_REDIRECT_URL", ""),
			},
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:     parseInt(getEnv("SMTP_PORT", "587")),
			SMTPUser:     getEnv("SMTP_USER", ""),
			SMTPPassword: getEnv("SMTP_PASSWORD", ""),
			SMTPFrom:     getEnv("SMTP_FROM", ""),
		},
		Logger: LoggerConfig{
			Level:  getEnv("LOG_LEVEL", "debug"),
			Format: getEnv("LOG_FORMAT", "json"),
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
	fmt.Printf("   JWT Secret: %s\n", maskSecret(cfg.JWT.Secret))
	fmt.Printf("   Log Level: %s\n", cfg.Logger.Level)
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
