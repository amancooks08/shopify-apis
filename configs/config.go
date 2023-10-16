package configs

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

var (
	appPort string
)

func Load() {
	viper.SetDefault("APP_PORT", "8002")

	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./..")
	viper.AddConfigPath("./../..")
	viper.ReadInConfig()
	viper.AutomaticEnv()
}

func AppPort() string {
	if appPort == "" {
		appPort = ReadEnvString("APP_PORT")
	}
	return appPort
}

func ReadEnvString(key string) string {
	checkIfSet(key)
	return viper.GetString(key)
}

func checkIfSet(key string) {
	if !viper.IsSet(key) {
		err := errors.New(fmt.Sprintf("Key %s is not set", key))
		panic(err)
	}
}