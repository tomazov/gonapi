package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppEnv             string
	ServerURL          string
	ServerReadTimeout  time.Duration
	JWTSecretKey       string
	JWTExpireMinutes   int

	DBConnectionString string
	DBMaxConnections   int
	DBMaxIdle          int
	DBMaxLifetime      int

	RabbitMQ RabbitMQConfig
	MemcachedURL string
	ClickHouse   ClickHouseConfig
}

type RabbitMQConfig struct {
	Host     string
	User     string
	Password string
	Port     int
	VHost    string
	WebPort  int
}

type ClickHouseConfig struct {
	Host     string
	User     string
	Password string
	Database string
}

var Cfg AppConfig

func Load() {
	_ = godotenv.Load()

	Cfg = AppConfig{
		AppEnv:             getEnv("APP_ENV", "development"),
		ServerURL:          getEnv("SERVER_URL", "0.0.0.0:5000"),
		ServerReadTimeout:  time.Second * time.Duration(getEnvAsInt("SERVER_READ_TIMEOUT", 30)),
		JWTSecretKey:       getEnv("JWT_SECRET_KEY", "secret"),
		JWTExpireMinutes:   getEnvAsInt("JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT", 15),

		DBConnectionString: getEnv("DB_SERVER_URL", ""),
		DBMaxConnections:   getEnvAsInt("DB_MAX_CONNECTIONS", 100),
		DBMaxIdle:          getEnvAsInt("DB_MAX_IDLE_CONNECTIONS", 10),
		DBMaxLifetime:      getEnvAsInt("DB_MAX_LIFETIME_CONNECTIONS", 2),

		RabbitMQ: RabbitMQConfig{
			Host:     getEnv("RABBITMQ_HOST", "localhost"),
			User:     getEnv("RABBITMQ_USER", "guest"),
			Password: getEnv("RABBITMQ_PASS", "guest"),
			Port:     getEnvAsInt("RABBITMQ_PORT", 5672),
			VHost:    getEnv("RABBITMQ_VHOST", "/"),
			WebPort:  getEnvAsInt("RABBITMQ_WEB_PORT", 15672),
		},

		MemcachedURL: getEnv("MEMCACHED_URL", "localhost:11211"),

		ClickHouse: ClickHouseConfig{
			Host:     getEnv("CLICKHOUSE_HOST", "localhost"),
			User:     getEnv("CLICKHOUSE_USER", "default"),
			Password: getEnv("CLICKHOUSE_PASSWORD", ""),
			Database: getEnv("CLICKHOUSE_DB", "default"),
		},
	}

	log.Printf("✅ Config loaded for ENV: %s on %s", Cfg.AppEnv, Cfg.ServerURL)
}

func getEnv(key string, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if val, ok := os.LookupEnv(key); ok {
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed
		}
	}
	return fallback
}
