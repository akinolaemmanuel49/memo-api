package internal

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
)

// Config holds configuration data passed via flags or dotenv
type Config struct {
	Env       string
	Port      int
	JWTSecret string

	Db struct {
		DSN string
	}

	Cloudinary struct {
		CloudName string
		APIKey    string
		APISecret string
	}
}

// Parse sets the fields of the Config to the data passed in via flags or dotenv
func (c *Config) Parse() {
	// app details
	flag.StringVar(&c.Env, "env", c.defaultEnv(), "Working Environment (development | staging | production)\nDotenv variable: ENV\n")
	flag.IntVar(&c.Port, "port", c.defaultPort(), "API Server Port\nDotenv variable: PORT\n")
	flag.StringVar(&c.JWTSecret, "jwt-secret", c.defaultJWTSecret(), "JWT Secret Key - Required\nDotenv variable: JWT_SECRET\n")

	// database details
	flag.StringVar(&c.Db.DSN, "db-dsn", c.defaultDbDSN(), "Postgres Database DSN - Required\nDotenv variable: DB_DSN\n")

	// cloudinary details
	flag.StringVar(&c.Cloudinary.CloudName, "cloudinary-cloud-name", c.defaultCloudinaryCloudName(), "Cloudinary Cloud Name\nDotenv variable: CLOUDINARY_CLOUD_NAME\n")
	flag.StringVar(&c.Cloudinary.APIKey, "cloudinary-api-key", c.defaultCloudinaryAPIKey(), "Cloudinary API Key\nDotenv variable: CLOUDINARY_API_KEY\n")
	flag.StringVar(&c.Cloudinary.APISecret, "cloudinary-api-secret", c.defaultCloudinaryAPISecret(), "Cloudinary API Secret\nDotenv variable: CLOUDINARY_API_SECRET\n")

	flag.Parse()
}

// Validate ensures required flags or environment variables are set.
func (c *Config) Validate() error {
	if c.JWTSecret == "" {
		return errors.New(validationMessage("jwt-secret", "JWT_SECRET"))
	}

	if c.Db.DSN == "" {
		return errors.New(validationMessage("db-dsn", "DB_DSN"))
	}

	return nil
}

func validationMessage(flag string, dotenv string) string {
	return fmt.Sprintf("the %q flag or %q dotenv variable is required", flag, dotenv)
}

func (c *Config) defaultEnv() string {
	const defaultEnv = "development"

	if env, exists := os.LookupEnv("ENV"); exists {
		return env
	}
	return defaultEnv
}

func (c *Config) defaultPort() int {
	const defaultPort = 5000

	if portEnv, exists := os.LookupEnv("PORT"); exists {
		port, err := strconv.Atoi(portEnv)
		if err == nil {
			return port
		}
	}
	return defaultPort
}

func (c *Config) defaultJWTSecret() string {
	const defaultSecret = ""

	if secret, exists := os.LookupEnv("JWT_SECRET"); exists {
		return secret
	}
	return defaultSecret
}

func (c *Config) defaultDbDSN() string {
	const defaultDSN = ""

	if dsn, exists := os.LookupEnv("DB_DSN"); exists {
		return dsn
	}
	return defaultDSN
}

func (c *Config) defaultCloudinaryCloudName() string {
	const defaultCloudinaryCloudName = ""

	if value, exists := os.LookupEnv("CLOUDINARY_CLOUD_NAME"); exists {
		return value
	}
	return defaultCloudinaryCloudName
}

func (c *Config) defaultCloudinaryAPIKey() string {
	const defaultCloudinaryAPIKey = ""

	if value, exists := os.LookupEnv("CLOUDINARY_API_KEY"); exists {
		return value
	}
	return defaultCloudinaryAPIKey
}

func (c *Config) defaultCloudinaryAPISecret() string {
	const defaultCloudinaryAPISecret = ""

	if value, exists := os.LookupEnv("CLOUDINARY_API_SECRET"); exists {
		return value
	}
	return defaultCloudinaryAPISecret
}
