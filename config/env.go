package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"strconv"
)

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env not found or error: %v", err)
	}
}

func Init() {
	loadEnv()
}

func getString(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	i, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return i
}

func getBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	b, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return b
}

type DatabaseConfig struct {
	Url string
}

func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Url: getString("DATABASE_URL", ""),
	}
}

type LogConfig struct {
	Level  int
	Format string
}

func NewLogConfig() *LogConfig {
	return &LogConfig{
		Level:  getInt("LOG_LEVEL", 0),
		Format: getString("LOG_FORMAT", "json"),
	}
}

type ImportConfig struct {
	ImportUrl string
	ImportKey string
}

func NewImportConfig() *ImportConfig {
	return &ImportConfig{
		ImportUrl: getString("IMPORT_URL", ""),
		ImportKey: getString("IMPORT_KEY", "2c70e12b7a0646f92279f427c7b38e7334d8e5389cff167a1dc30e73f826b683"),
	}
}

func GetAppPort() string {
	return getString("APP_PORT", "8080")
}
