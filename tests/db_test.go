package tests

import (
	"Avito-backend-trainee-assignment-winter-2025/config"
	"Avito-backend-trainee-assignment-winter-2025/internal/database"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestConnectDB(t *testing.T) {
	cfg := &config.Config{
		DBHost: "localhost",
		DBUser: "avito",
		DBPass: "secret",
		DBName: "avito_shop",
		DBPort: "5432",
	}

	log := zerolog.New(nil)

	db, err := database.ConnectDB(cfg, &log)

	assert.NoError(t, err)
	assert.NotNil(t, db)

	database.CloseDB(db, &log)
}

func TestConnectDB_InvalidConfig(t *testing.T) {
	cfg := &config.Config{
		DBHost: "invalid_host",
		DBUser: "invalid_user",
		DBPass: "invalid_password",
		DBName: "invalid_db",
		DBPort: "invalid_port",
	}

	log := zerolog.New(nil)

	db, err := database.ConnectDB(cfg, &log)

	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestCloseDB(t *testing.T) {
	cfg := &config.Config{
		DBHost: "localhost",
		DBUser: "avito",
		DBPass: "secret",
		DBName: "avito_shop",
		DBPort: "5432",
	}

	log := zerolog.New(nil)

	db, err := database.ConnectDB(cfg, &log)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	database.CloseDB(db, &log)

	dbSQL, err := db.DB()
	assert.NoError(t, err)
	assert.NoError(t, dbSQL.Close())
}
