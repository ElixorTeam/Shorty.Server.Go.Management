package config

import (
	"Shorty.Server.Go.Management/internal/constants"
	"github.com/spf13/viper"
)

var AppConfig Config

type Config struct {
	Port int `mapstructure:"PORT"`

	DBDriver   string `mapstructure:"DB_DRIVER"`
	DBDatabase string `mapstructure:"DB_DATABASE"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBUser     string `mapstructure:"DB_USER"`
	DBSchema   string `mapstructure:"DB_SCHEMA"`

	AUTHJwkPublicUri      string `mapstructure:"AUTH_JWK_SECRET_URI"`
	AUTHRealm             string `mapstructure:"AUTH_REALM"`
	AUTHClient            string `mapstructure:"AUTH_CLIENT"`
	AUTHUserInfoEndpoint  string `mapstructure:"AUTH_USER_INFO_ENDPOINT"`
	AUTHRefreshJwkTimeout int    `mapstructure:"AUTH_REFRESH_JWK_TIMEOUT"`

	RedisAddress  string `mapstructure:"REDIS_ADDRESS"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
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
	err = viper.ReadInConfig()
	if err != nil {
		return constants.ErrLoadConfig
	}

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		return constants.ErrParseConfig
	}

	if AppConfig.Port == 0 || AppConfig.DBHost == "" || AppConfig.DBPassword == "" || AppConfig.DBUser == "" || AppConfig.DBSchema == "" || AppConfig.DBDatabase == "" {
		return constants.ErrEmptyVar
	}

	return nil
}

func init() {
	err := InitializeAppConfig()
	if err != nil {
		panic(err)
	}
}
