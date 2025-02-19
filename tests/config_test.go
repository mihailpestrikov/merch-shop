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

	config, err := config2.LoadConfig("./testconfig")

	assert.NoError(t, err)

	assert.Equal(t, "TestApp", config.AppName)
	assert.Equal(t, "localhost", config.AppHost)
	assert.Equal(t, "8050", config.AppPort)
	assert.Equal(t, "localhost", config.DBHost)
	assert.Equal(t, "root", config.DBUser)
	assert.Equal(t, "password", config.DBPass)
	assert.Equal(t, "test_db", config.DBName)
	assert.Equal(t, "5431", config.DBPort)
	assert.Equal(t, "disable", config.DBSslMode)
	assert.Equal(t, "debug", config.LogLevel)
	assert.Equal(t, "json", config.LogFormat)
	assert.Equal(t, "secret", config.SecretKey)
	assert.Equal(t, 10*time.Second, config.ShutdownTimeout)
	assert.Equal(t, 10, config.MaxIdleConns)
	assert.Equal(t, 50, config.MaxOpenConns)
	assert.Equal(t, 60*time.Minute, config.ConnMaxLifetime)
}
