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
	DBSource             string        `mapstructure:"DB_SOURCE"`
	Port                 string        `mapstructure:"PORT"`
	PortGatewayGrpc      string        `mapstructure:"PORT_GATEWAY_GRPC"`
	GRPCPort             string        `mapstructure:"GRPC_PORT"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	Environment          string        `mapstructure:"ENVIROMENT"`
	RedisPort            string        `mapstructure:"REDIS_PORT"`
	RedisHost            string        `mapstructure:"REDIS_HOST"`
	RedisDB              string        `mapstructure:"REDIS_DB"`
	MailHost             string        `mapstructure:"MAIL_HOST"`
	MailPort             int           `mapstructure:"MAIL_PORT"`
	MailUser             string        `mapstructure:"MAIL_USER"`
	MailPassword         string        `mapstructure:"MAIL_PASSWORD"`
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
		_ = os.Setenv("DB_SOURCE", viper.GetString("DB_SOURCE"))
		_ = os.Setenv("PORT", viper.GetString("PORT"))
		_ = os.Setenv("PORT_GATEWAY_GRPC", viper.GetString("PORT_GATEWAY_GRPC"))
		_ = os.Setenv("GRPC_PORT", viper.GetString("GRPC_PORT"))
		_ = os.Setenv("TOKEN_SYMMETRIC_KEY", viper.GetString("TOKEN_SYMMETRIC_KEY"))
		_ = os.Setenv("ACCESS_TOKEN_DURATION", viper.GetString("ACCESS_TOKEN_DURATION"))
		_ = os.Setenv("REFRESH_TOKEN_DURATION", viper.GetString("REFRESH_TOKEN_DURATION"))
		_ = os.Setenv("ENVIROMENT", viper.GetString("ENVIROMENT"))
		_ = os.Setenv("REDIS_PORT", viper.GetString("REDIS_PORT"))
		_ = os.Setenv("REDIS_HOST", viper.GetString("REDIS_HOST"))
		_ = os.Setenv("REDIS_DB", viper.GetString("REDIS_DB"))
		_ = os.Setenv("MAIL_HOST", viper.GetString("MAIL_HOST"))
		_ = os.Setenv("MAIL_PORT", viper.GetString("MAIL_PORT"))
		_ = os.Setenv("MAIL_USER", viper.GetString("MAIL_USER"))
		_ = os.Setenv("MAIL_PASSWORD", viper.GetString("MAIL_PASSWORD"))

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
		viper.BindEnv("DB_SOURCE")
		viper.BindEnv("PORT")
		viper.BindEnv("PORT_GATEWAY_GRPC")
		viper.BindEnv("GRPC_PORT")
		viper.BindEnv("TOKEN_SYMMETRIC_KEY")
		viper.BindEnv("ACCESS_TOKEN_DURATION")
		viper.BindEnv("REFRESH_TOKEN_DURATION")
		viper.BindEnv("ENVIROMENT")
		viper.BindEnv("REDIS_PORT")
		viper.BindEnv("REDIS_HOST")
		viper.BindEnv("REDIS_DB")
		viper.BindEnv("MAIL_HOST")
		viper.BindEnv("MAIL_PORT")
		viper.BindEnv("MAIL_USER")
		viper.BindEnv("MAIL_PASSWORD")

		// Unmarshal the config into the struct
		err = viper.Unmarshal(&config)
		if err != nil {
			return config, err
		}

		return config, nil
	}
}
