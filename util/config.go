package util

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost               string        `mapstructure:"DB_HOST"`
	DBPort               string        `mapstructure:"DB_PORT"`
	DBUser               string        `mapstructure:"DB_USER"`
	DBPassword           string        `mapstructure:"DB_PASSWORD"`
	DBName               string        `mapstructure:"DB_NAME"`
	DBSSLMode            string        `mapstructure:"DB_SSLMODE"`
	Port                 string        `mapstructure:"PORT"`
	PortGatewayGrpc      string        `mapstructure:"PORT_GATEWAY_GRPC"`
	GRPCPort             string        `mapstructure:"GRPC_PORT"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
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

		// set OS environment variables
		_ = os.Setenv("DB_HOST", viper.GetString("DB_HOST"))
		_ = os.Setenv("DB_PORT", viper.GetString("DB_PORT"))
		_ = os.Setenv("DB_USER", viper.GetString("DB_USER"))
		_ = os.Setenv("DB_PASSWORD", viper.GetString("DB_PASSWORD"))
		_ = os.Setenv("DB_NAME", viper.GetString("DB_NAME"))
		_ = os.Setenv("DB_SSLMODE", viper.GetString("DB_SSLMODE"))
		_ = os.Setenv("PORT", viper.GetString("PORT"))
		_ = os.Setenv("PORT_GATEWAY_GRPC", viper.GetString("PORT_GATEWAY_GRPC"))
		_ = os.Setenv("GRPC_PORT", viper.GetString("GRPC_PORT"))
		_ = os.Setenv("TOKEN_SYMMETRIC_KEY", viper.GetString("TOKEN_SYMMETRIC_KEY"))
		_ = os.Setenv("ACCESS_TOKEN_DURATION", viper.GetString("ACCESS_TOKEN_DURATION"))
		_ = os.Setenv("REFRESH_TOKEN_DURATION", viper.GetString("REFRESH_TOKEN_DURATION"))

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
		viper.BindEnv("PORT_GATEWAY_GRPC")
		viper.BindEnv("GRPC_PORT")
		viper.BindEnv("TOKEN_SYMMETRIC_KEY")
		viper.BindEnv("ACCESS_TOKEN_DURATION")
		viper.BindEnv("REFRESH_TOKEN_DURATION")

		// Unmarshal the config into the struct
		err = viper.Unmarshal(&config)
		if err != nil {
			return config, err
		}

		return config, nil
	}
}
