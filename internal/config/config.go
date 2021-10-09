package config

import (
	"github.com/joho/godotenv"
	"github.com/online-shop/pkg/log"
	"os"
	"strconv"
)

const (
	defaultServerPort         = 8080
	defaultJWTExpirationHours = 72
)

// Config represents an application configuration.
type Config struct {
	// the server port. Defaults to 8080
	ServerPort int `env:"SERVER_PORT"`
	// the data source name (DSN) for connecting to the database. required.
	DSN string `env:"DSN,secret"`
	// JWT signing key. required.
	JWTSigningKey string `env:"JWT_SIGNING_KEY,secret"`
	// JWT expiration in hours. Defaults to 72 hours (3 days)
	JWTExpiration int `env:"JWT_EXPIRATION"`
}

func Load(logger log.Logger) (*Config, error) {
	// default config
	c := Config{
		ServerPort:    defaultServerPort,
		JWTExpiration: defaultJWTExpirationHours,
	}

	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
	}

	c.ServerPort = getEnvAsInt("SERVER_PORT", defaultServerPort)
	c.DSN = os.Getenv("DSN")
	//secretKey := os.Getenv("SECRET_KEY")

	return &c, err
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := os.Getenv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}
