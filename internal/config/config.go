package config

import "github.com/spf13/viper"

type Config_PG struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSLMODE"`
}

func Load_Config_PG() (*Config_PG, error) {
	viper.SetConfigFile(".env") // Указываем файл с конфигом
	viper.AutomaticEnv()        // читает переменные из окружения(переопределяя переменные из .env)

	// Читаем и загружаем файл конфига
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config_PG
	// Преобразуем данные которые получили в нашу структуру(Config_PG)
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
