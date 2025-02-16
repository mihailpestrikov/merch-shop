package service

import (
	"Avito-backend-trainee-assignment-winter-2025/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(username, password string) error
	Login(username, password string) (string, error)
}

type AuthServiceImpl struct {
	db          *gorm.DB
	log         *zerolog.Logger
	userService UserService
	cfg         *config.Config
}

func NewAuthService(db *gorm.DB, log *zerolog.Logger, userService UserService, cfg *config.Config) *AuthServiceImpl {
	return &AuthServiceImpl{db: db, log: log, userService: userService, cfg: cfg}
}

func (s *AuthServiceImpl) Register(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to hash password")

		return err
	}
	password = string(hashedPassword)

	if _, err := s.userService.CreateUser(username, password); err != nil {
		return err
	}

	return nil
}

func (s *AuthServiceImpl) Login(username, password string) (string, error) {
	user, err := s.userService.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	curTime := time.Now().Add(time.Hour * 24).Unix()
	s.log.Info().Int64("curTime", curTime).Msg("Expire date")
	s.log.Info().Str("username", user.Username).Msg("Username")
	s.log.Info().Str("user_id", user.ID.String()).Msg("User id")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": username,
		"exp":      curTime,
	})

	tokenString, err := token.SignedString([]byte(s.cfg.SecretKey))
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to generate token")

		return "", err
	}

	return tokenString, nil
}
