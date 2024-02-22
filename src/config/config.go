package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	MaxConnections int
	Port           int
}

func GetConfig() (*Config, error) {
	err := godotenv.Load()

	if err != nil {
		return nil, err
	}

	return &Config{
		MaxConnections: readEnvNumber("NOSTRESSGO_MAX_CONN", 5),
		Port:           readEnvNumber("NOSTRESSGO_PORT", 3000),
	}, nil
}

func (c *Config) GetPostgresConnString() (string, error) {
	username, err := readEnvVar("PGUSER")
	if err != nil {
		return "", err
	}

	password, err := readEnvVar("PGPASSWORD")
	if err != nil {
		return "", err
	}

	host, err := readEnvVar("PGHOST")
	if err != nil {
		return "", err
	}

	port, err := readEnvVar("PGPORT")
	if err != nil {
		return "", err
	}

	dbname, err := readEnvVar("PGDATABASE")
	if err != nil {
		return "", err
	}

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, dbname), nil
}

func readEnvNumber(key string, defaultValue int) int {
	value, err := readEnvVar(key)
	if err != nil {
		return defaultValue
	}

	i64, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return defaultValue
	}

	return int(i64)
}

func readEnvVar(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("missing env var %s", key)
	}
	return value, nil
}
