package tests

import (
	"Avito-backend-trainee-assignment-winter-2025/config"
	"Avito-backend-trainee-assignment-winter-2025/internal/logger"
	"bytes"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestInitLogger(t *testing.T) {
	t.Run("Default Level", func(t *testing.T) {
		cfg := &config.Config{
			LogLevel:  "invalid_level",
			LogFormat: "json",
		}

		log := logger.InitLogger(cfg)
		assert.NotNil(t, log)
		assert.Equal(t, zerolog.InfoLevel, log.GetLevel())
	})

	t.Run("Text Format", func(t *testing.T) {
		cfg := &config.Config{
			LogLevel:  "debug",
			LogFormat: "text",
		}

		log := logger.InitLogger(cfg)
		assert.NotNil(t, log)
		assert.Equal(t, zerolog.DebugLevel, log.GetLevel())

		output := bytes.NewBuffer(nil)
		consoleWriter := zerolog.ConsoleWriter{Out: output}
		textLogger := zerolog.New(consoleWriter).Level(zerolog.DebugLevel).With().Timestamp().Logger()
		textLogger.Info().Msg("Test message")

		assert.Contains(t, output.String(), "Test message")
	})

	t.Run("JSON Format", func(t *testing.T) {
		cfg := &config.Config{
			LogLevel:  "error",
			LogFormat: "json",
		}

		log := logger.InitLogger(cfg)
		assert.NotNil(t, log)
		assert.Equal(t, zerolog.ErrorLevel, log.GetLevel())

		output := bytes.NewBuffer(nil)
		jsonLogger := zerolog.New(output).Level(zerolog.ErrorLevel).With().Timestamp().Logger()
		jsonLogger.Error().Msg("Test message")

		assert.Contains(t, output.String(), "Test message")
	})
}
