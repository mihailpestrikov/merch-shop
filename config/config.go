package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	AppName   string
	AppEnv    string
	AppHost   string
	AppPort   string
	DBHost    string
	DBUser    string
	DBPass    string
	DBName    string
	DBPort    string
	LogLevel  string
	LogFormat string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file, %s", err)
	}

	viper.AutomaticEnv()
	env := viper.GetString("ENV")

	if env == "" {
		env = "local"
	}

	viper.SetConfigName("config." + env)
	if err := viper.MergeInConfig(); err != nil {
		return nil, fmt.Errorf("config file for env: %s not found: %e", env, err)
	}

	AppName := viper.GetString("app.name")
	AppHost := viper.GetString("app.host")
	AppPort := viper.GetString("app.port")
	DBHost := viper.GetString("db.host")
	DBUser := viper.GetString("db.user")
	DBPass := viper.GetString("db.password")
	DBName := viper.GetString("db.name")
	DBPort := viper.GetString("db.port")
	DBSslMode := viper.GetString("db.ssl-mode")
	LogLevel := viper.GetString("log.level")
	LogFormat := viper.GetString("log.format")

	if AppName == "" || AppHost == "" || AppPort == "" {
		return nil, fmt.Errorf("one or more app configuration fields are empty")
	}

	if DBHost == "" || DBUser == "" || DBPass == "" || DBName == "" || DBPort == "" || DBSslMode == "" {
		return nil, fmt.Errorf("one or more database configuration fields are empty")
	}

	if LogLevel == "" || LogFormat == "" {
		return nil, fmt.Errorf("one or more logger configuration fields are empty")
	}

	config := Config{
		AppName:   AppName,
		AppHost:   AppHost,
		AppPort:   AppPort,
		AppEnv:    env,
		DBHost:    DBHost,
		DBUser:    DBUser,
		DBPass:    DBPass,
		DBName:    DBName,
		DBPort:    DBPort,
		LogLevel:  LogLevel,
		LogFormat: LogFormat,
	}

	return &config, nil
}
