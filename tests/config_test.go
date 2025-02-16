package tests

import (
	config2 "Avito-backend-trainee-assignment-winter-2025/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLoadConfig_Success(t *testing.T) {
	viper.Reset()

	viper.Set("app.name", "TestApp")
	viper.Set("app.host", "localhost")
	viper.Set("app.port", "8080")
	viper.Set("db.host", "localhost")
	viper.Set("db.user", "root")
	viper.Set("db.password", "password")
	viper.Set("db.name", "test_db")
	viper.Set("db.port", "5432")
	viper.Set("db.ssl-mode", "disable")
	viper.Set("log.level", "debug")
	viper.Set("log.format", "json")
	viper.Set("auth.secret-key", "secret")
	viper.Set("app.shutdown-timeout", "10")

	config, err := config2.LoadConfig("../config")

	assert.NoError(t, err)

	assert.Equal(t, "TestApp", config.AppName)
	assert.Equal(t, "localhost", config.AppHost)
	assert.Equal(t, "8080", config.AppPort)
	assert.Equal(t, "localhost", config.DBHost)
	assert.Equal(t, "root", config.DBUser)
	assert.Equal(t, "password", config.DBPass)
	assert.Equal(t, "test_db", config.DBName)
	assert.Equal(t, "5432", config.DBPort)
	assert.Equal(t, "disable", config.DBSslMode)
	assert.Equal(t, "debug", config.LogLevel)
	assert.Equal(t, "json", config.LogFormat)
	assert.Equal(t, "secret", config.SecretKey)
	assert.Equal(t, 10*time.Second, config.ShutdownTimeout)
}
