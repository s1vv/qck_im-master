package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config - структура для хранения настроек
type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPass        string
	DBName        string
	SMTPHost      string
	SMTPPort      int
	SMTPUser      string
	SMTPPass      string
	JWTSecret     string
	JWTExpiration int
	LogLevel      string
	ServerPort    string
	AppEnv        string
	BaseURL       string
}

var cfg *Config

// LoadConfig загружает переменные из `.env`
func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		println("No .env file found, using system environment variables", "error", err)
	}

	cfg = &Config{
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "3306"),
		DBUser:        getEnv("DB_USER", "root"),
		DBPass:        getEnv("DB_PASS", "root"),
		DBName:        getEnv("DB_NAME", "qck"),
		SMTPHost:      getEnv("SMTP_HOST", "smtp.yandex.ru"),
		SMTPPort:      getEnvAsInt("SMTP_PORT", 587),
		SMTPUser:      getEnv("SMTP_USER", ""),
		SMTPPass:      getEnv("SMTP_PASS", ""),
		JWTSecret:     getEnv("JWT_SECRET", ""),
		JWTExpiration: getEnvAsInt("JWT_EXPIRATION", 86400),
		LogLevel:      getEnv("LOG_LEVEL", "Debug"),
		ServerPort:    getEnv("SERVER_PORT", "8081"),
		AppEnv:        getEnv("APP_ENV", "development"),
		BaseURL:       getEnv("BASE_URL", "http://localhost:80"),
	}
}

// GetConfig возвращает загруженную конфигурацию
func GetConfig() *Config {
	if cfg == nil {
		LoadConfig()
	}
	return cfg
}

// getEnv получает строковое значение из переменной окружения
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt получает числовое значение из переменной окружения
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
