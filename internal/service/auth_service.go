package service

import (
	"Avito-backend-trainee-assignment-winter-2025/config"
	"Avito-backend-trainee-assignment-winter-2025/internal/models"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/scrypt"
)

const (
	scryptN    = 8192 // Complexity (higher = slower, more secure)
	scryptR    = 4    // Memory block size (higher = more memory)
	scryptP    = 1    // Parallelization factor (higher = more threads)
	saltLength = 16   // Salt length in bytes
	hashLength = 32   // Length of the resulting hash in bytes
)

type AuthService interface {
	LoginOrRegister(username, password string) (string, error)
}

type AuthServiceImpl struct {
	log         *zerolog.Logger
	userService UserService
	cfg         *config.Config
}

func NewAuthService(log *zerolog.Logger, userService UserService, cfg *config.Config) *AuthServiceImpl {
	return &AuthServiceImpl{log: log, userService: userService, cfg: cfg}
}

func hashPassword(password string) (string, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash, err := scrypt.Key([]byte(password), salt, scryptN, scryptR, scryptP, hashLength)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(append(salt, hash...)), nil
}

func verifyPassword(password, stored string) bool {
	data, err := base64.StdEncoding.DecodeString(stored)
	if err != nil || len(data) < saltLength {
		return false
	}

	salt, hash := data[:saltLength], data[saltLength:]
	newHash, err := scrypt.Key([]byte(password), salt, scryptN, scryptR, scryptP, hashLength)
	if err != nil {
		return false
	}

	return string(newHash) == string(hash)
}

func (s *AuthServiceImpl) LoginOrRegister(username, password string) (string, error) {
	user, err := s.userService.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			hashedPassword, err := hashPassword(password)
			if err != nil {
				s.log.Error().Err(err).Msg("Failed to hash password")
				return "", err
			}

			user, err = s.userService.CreateUser(username, hashedPassword)
			if err != nil {
				s.log.Error().Err(err).Msg("Failed to create user")
				return "", err
			}
		} else {
			return "", err
		}
	}

	if !verifyPassword(password, user.Password) {
		return "", errors.New("invalid password")
	}

	expTime := time.Now().Add(24 * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID.String(),
		"username": username,
		"exp":      expTime,
	})

	tokenString, err := token.SignedString([]byte(s.cfg.SecretKey))
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to generate token")
		return "", err
	}

	return tokenString, nil
}
