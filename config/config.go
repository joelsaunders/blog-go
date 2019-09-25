package config

import (
	_ "github.com/lib/pq"
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
	viper.SetConfigName("bilbo.config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()

	if err != nil {
		return Constants{}, err
	}

	viper.SetDefault("PORT", "8000")
	viper.SetDefault("JwtSecret", "not so secret")
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
