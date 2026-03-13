package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	JWT      JWTConfig
	Database DatabaseConfig
	Upload   UploadConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type JWTConfig struct {
	Secret              string
	ExpireHours         int
	RefreshTokenExpireDays int
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type UploadConfig struct {
	MaxFileSize int64
	Path        string
}

func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("SERVER_MODE", "debug"),
		},
		JWT: JWTConfig{
			Secret:                 getEnv("JWT_SECRET", "nusvakspps-secret-key"),
			ExpireHours:            getEnvAsInt("JWT_EXPIRE_HOURS", 24),
			RefreshTokenExpireDays: getEnvAsInt("REFRESH_TOKEN_EXPIRE_DAYS", 7),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "nusvakspps"),
		},
		Upload: UploadConfig{
			MaxFileSize: int64(getEnvAsInt("MAX_FILE_SIZE", 2097152)), // 2MB
			Path:        getEnv("UPLOAD_PATH", "./uploads"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
