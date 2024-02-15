package config

import (
	"Shorty.Server.Go.Management/internal/constants"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"log"
)

var AppConfig Config

type Config struct {
	Port  int  `mapstructure:"PORT"`
	Debug bool `mapstructure:"DEBUG"`

	DBDriver   string `mapstructure:"DB_DRIVER" validate:"required"`
	DBName     string `mapstructure:"DB_NAME" validate:"required"`
	DBHost     string `mapstructure:"DB_HOST" validate:"required"`
	DBPort     int16  `mapstructure:"DB_PORT" validate:"required"`
	DBPassword string `mapstructure:"DB_PASSWORD" validate:"required"`
	DBUser     string `mapstructure:"DB_USER" validate:"required"`
	DBSchema   string `mapstructure:"DB_SCHEMA" validate:"required"`

	AUTHJwkPublicUri      string `mapstructure:"AUTH_JWK_SECRET_URI" validate:"required"`
	AUTHRealm             string `mapstructure:"AUTH_REALM" validate:"required"`
	AUTHClient            string `mapstructure:"AUTH_CLIENT" validate:"required"`
	AUTHUserInfoEndpoint  string `mapstructure:"AUTH_USER_INFO_ENDPOINT" validate:"required"`
	AUTHRefreshJwkTimeout int    `mapstructure:"AUTH_REFRESH_JWK_TIMEOUT" validate:"required"`

	RedisAddress  string `mapstructure:"REDIS_ADDRESS" validate:"required"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD" validate:"required"`
	RedisDB       int    `mapstructure:"REDIS_DB"`
}

func InitializeAppConfig() (err error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("internal/config")
	viper.AddConfigPath("/")
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()

	viper.SetDefault("PORT", 2000)
	viper.SetDefault("DEBUG", false)

	err = viper.ReadInConfig()
	if err != nil {
		return constants.ErrLoadConfig
	}

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		return constants.ErrParseConfig
	}

	validate := validator.New()
	if err := validate.Struct(&AppConfig); err != nil {
		log.Fatalf("Missing required attributes %v\n", err)
	}

	return nil
}

func init() {
	err := InitializeAppConfig()
	if err != nil {
		panic(err)
	}
}
