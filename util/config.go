package util

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost              string        `mapstructure:"DB_HOST"`
	DBPort              string        `mapstructure:"DB_PORT"`
	DBUser              string        `mapstructure:"DB_USER"`
	DBPassword          string        `mapstructure:"DB_PASSWORD"`
	DBName              string        `mapstructure:"DB_NAME"`
	DBSSLMode           string        `mapstructure:"DB_SSLMODE"`
	Port                string        `mapstructure:"PORT"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	if os.Getenv("DB_USER") == "" || os.Getenv("DB_PASSWORD") == "" || os.Getenv("DB_NAME") == "" {
		viper.AddConfigPath(path)
		viper.SetConfigName("app")
		viper.SetConfigType("env")

		viper.AutomaticEnv()

		err = viper.ReadInConfig()
		if err != nil {
			return
		}

		err = viper.Unmarshal(&config)
		return config, err
	} else {
		viper.AutomaticEnv()

		// Bind environment variables to struct fields
		viper.BindEnv("DB_HOST")
		viper.BindEnv("DB_PORT")
		viper.BindEnv("DB_USER")
		viper.BindEnv("DB_PASSWORD")
		viper.BindEnv("DB_NAME")
		viper.BindEnv("DB_SSLMODE")
		viper.BindEnv("PORT")
		viper.BindEnv("TOKEN_SYMMETRIC_KEY")
		viper.BindEnv("ACCESS_TOKEN_DURATION")

		// Unmarshal the config into the struct
		err = viper.Unmarshal(&config)
		if err != nil {
			return config, err
		}

		return config, nil
	}
}
