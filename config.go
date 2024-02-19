package nostressgo

import (
	"os"
	"strconv"
)

type Config struct {
	maxConnections int
}

func GetConfig() *Config {
	return &Config{
		maxConnections: readEnvNumber("NOSTRESSGO_MAX_CONN", 5),
	}
}

func readEnvNumber(key string, defaultValue int) int {
	value := readEnvVar(key)
	i64, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return defaultValue
	}
	return int(i64)
}

func readEnvVar(key string) string {
	return os.Getenv(key)
}
