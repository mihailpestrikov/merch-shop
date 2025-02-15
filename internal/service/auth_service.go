package service

import (
	"Avito-backend-trainee-assignment-winter-2025/config"
	"Avito-backend-trainee-assignment-winter-2025/internal/models"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type AuthService interface {
	Register(user models.User) error
	Login(username, password string) (string, error)
	ValidateToken(token string) (uuid.UUID, error)
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

func (s *AuthServiceImpl) Register(user models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to hash password")
		return err
	}
	user.Password = string(hashedPassword)

	if _, err := s.userService.CreateUser(user.Username, user.Password); err != nil {
		s.log.Error().Err(err).Msg("Failed to create user")
		return err
	}

	s.log.Info().Msg("User created")
	return nil
}

func (s *AuthServiceImpl) Login(username, password string) (string, error) {
	var user models.User
	if _, err := s.userService.GetUserByUsername(username); err != nil {
		s.log.Error().Err(err).Msg("Failed to find user")
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.cfg.SecretKey))
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to generate token")
		return "", err
	}

	return tokenString, nil
}

func (s *AuthServiceImpl) ValidateToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.cfg.SecretKey), nil
	})
	if err != nil {
		s.log.Error().Err(err).Msg("Could not parse token")
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, err := uuid.Parse(claims["user_id"].(string))
		if err != nil {
			return uuid.Nil, err
		}
		s.log.Info().Err(err).Msg("Token is valid")
		return userID, nil
	}

	s.log.Info().Err(err).Msg("Token is invalid")
	return uuid.Nil, errors.New("invalid token")
}
