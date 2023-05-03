package config

import (
	"fmt"
	"os"
)

// Структура переменных окружения
type AppEnv struct {
	LogLvl            string
	MinioRootUser     string
	MinioRootPassword string
}

// Формирование структуры перемемнных окружения
func GetEnvironment() (env AppEnv, err error) {
	env = AppEnv{
		LogLvl:            getEnv("LOG_LEVEL", "debug"),
		MinioRootUser:     getEnv("MINIO_ROOT_USER", ""),
		MinioRootPassword: getEnv("MINIO_ROOT_PASSWORD", ""),
	}

	if env.MinioRootUser == "" || env.MinioRootPassword == "" {
		return env, fmt.Errorf("incorrect environment params")
	}

	return env, nil
}

// Получение переменной окружения
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
