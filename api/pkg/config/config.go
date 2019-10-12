package config

import (
	"fmt"
	"strings"

	viper "github.com/spf13/viper"
)

type Constants struct {
	PORT     string
	Postgres struct {
		URL        string
		DBNAME     string
		DBPORT     string
		DBUSER     string
		DBPASSWORD string
	}
	JWTSecret []byte
}

type Config struct {
	Constants
}

func initViper() (Constants, error) {
	viper.SetConfigName("blog.config")
	viper.AddConfigPath(".")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()

	if err != nil {
		return Constants{}, err
	}

	viper.SetDefault("PORT", "8000")
	viper.Set("JWTSecret", []byte(viper.GetString("JWTSecret")))
	fmt.Println(viper.GetString("JWTSecret"))
	var constants Constants
	err = viper.Unmarshal(&constants)
	return constants, err
}

func NewConfig() (*Config, error) {
	config := Config{}
	constants, err := initViper()
	config.Constants = constants

	if err != nil {
		return &config, err
	}

	return &config, nil
}
