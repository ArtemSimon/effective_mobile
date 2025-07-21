package config

import (
	"effective_mobile/pkg/logger_module"

	"github.com/spf13/viper"
)

type Config_PG struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSLMODE"`
	Http_Port  string `mapstructure:"HTTP_PORT"`
}

func Load_Config_PG(logger *logger_module.Logger) (*Config_PG, error) {
	viper.SetConfigFile(".env") // Указываем файл с конфигом
	viper.AutomaticEnv()        // читает переменные из окружения(переопределяя переменные из .env)

	// Читаем и загружаем файл конфига
	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal("Failed to read config", "error", err)
	}

	var config Config_PG
	// Преобразуем данные которые получили в нашу структуру(Config_PG)
	if err := viper.Unmarshal(&config); err != nil {
		logger.Fatal("Failed unmarshal config struct", "error", err)
	}

	return &config, nil
}
