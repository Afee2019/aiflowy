package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	Log       LogConfig       `mapstructure:"log"`
	Snowflake SnowflakeConfig `mapstructure:"snowflake"`
	Storage   StorageConfig   `mapstructure:"storage"`
	Security  SecurityConfig  `mapstructure:"security"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Database        string `mapstructure:"database"`
	Charset         string `mapstructure:"charset"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"`
	Issuer string `mapstructure:"issuer"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
}

type SnowflakeConfig struct {
	WorkerID     int64 `mapstructure:"worker_id"`
	DatacenterID int64 `mapstructure:"datacenter_id"`
}

type StorageConfig struct {
	LocalRoot string `mapstructure:"local_root"`
}

type SecurityConfig struct {
	ApiKeyMasterKey string `mapstructure:"api_key_master_key"` // Bot API Key 加密主密钥 (32字节)
}

// DSN returns the database connection string
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.Charset,
	)
}

// Addr returns redis address
func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

var cfg *Config
var env string

// Load loads configuration from file
func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// Support environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Set defaults for snowflake
	if cfg.Snowflake.WorkerID == 0 {
		cfg.Snowflake.WorkerID = 1
	}
	if cfg.Snowflake.DatacenterID == 0 {
		cfg.Snowflake.DatacenterID = 1
	}

	// Set defaults for JWT
	if cfg.JWT.Secret == "" {
		cfg.JWT.Secret = "aiflowy-go-secret-key-change-in-production"
	}
	if cfg.JWT.Expire == 0 {
		cfg.JWT.Expire = 86400 // 24 hours
	}
	if cfg.JWT.Issuer == "" {
		cfg.JWT.Issuer = "aiflowy-go"
	}

	// Determine environment
	env = os.Getenv("GO_ENV")
	if env == "" {
		if cfg.Server.Mode == "release" {
			env = "production"
		} else {
			env = "development"
		}
	}

	return cfg, nil
}

// Get returns the global configuration
func Get() *Config {
	return cfg
}

// GetConfig returns the global configuration (alias for Get)
func GetConfig() *Config {
	return cfg
}

// GetEnv returns the current environment
func GetEnv() string {
	return env
}

// IsProduction returns true if running in production mode
func IsProduction() bool {
	return env == "production" || cfg.Server.Mode == "release"
}

// IsDevelopment returns true if running in development mode
func IsDevelopment() bool {
	return env == "development" || cfg.Server.Mode == "debug"
}
