package config

import (
	"github.com/spf13/viper"
)

// Конфиг сервера
type ServerConfig struct {
	Port string `mapstructure:"port"`
}

// MinIO конфиг
type MinioConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

// Конфиг приложения
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Minio  MinioConfig  `mapstructure:"minio"`
}

var vp *viper.Viper

// Загрузка конфига
func LoadConfig() (Config, error) {
	// Создание viper объекта конфига
	vp = viper.New()

	// Создание объекта конфига
	var config Config

	// Параметры конфига
	vp.SetConfigName("config")
	vp.SetConfigType("json")
	vp.AddConfigPath("config")

	// Чтение из config/config.json
	err := vp.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	// Десериализация json конфига
	err = vp.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
