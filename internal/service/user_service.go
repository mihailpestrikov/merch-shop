package service

import (
	"Avito-backend-trainee-assignment-winter-2025/internal/models"
	repository2 "Avito-backend-trainee-assignment-winter-2025/internal/repository"
	"errors"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"time"
)

type UserService interface {
	TransferCoins(fromUserID, toUserID uuid.UUID, amount int) error
	GetBalance(userID uuid.UUID) (int, error)
	PurchaseMerch(userID uuid.UUID, merchItemName string) error
	GetPurchasedItems(userID uuid.UUID) ([]models.MerchItem, error)
	GetTransactionHistory(userID uuid.UUID) ([]models.Transaction, error)
	CreateUser(username, passwordHash string) (*models.User, error)
	GetUserByID(userID uuid.UUID) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(userID uuid.UUID) error
}

type UserServiceImpl struct {
	userRepo        repository2.UserRepository
	transactionRepo repository2.TransactionRepository
	merchService    MerchService
	log             *zerolog.Logger
}

func NewUserService(userRepo repository2.UserRepository, transactionRepo repository2.TransactionRepository, merchService MerchService, log *zerolog.Logger) *UserServiceImpl {
	return &UserServiceImpl{userRepo: userRepo, transactionRepo: transactionRepo, merchService: merchService, log: log}
}

func (s *UserServiceImpl) GetBalance(userID uuid.UUID) (int, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		s.log.Error().Err(err).Msgf("Failed to get balance of user with ID: %s", userID)
		return 0, err
	}
	return user.Balance, nil
}

func (s *UserServiceImpl) TransferCoins(fromUserID, toUserID uuid.UUID, amount int) error {
	fromUser, err := s.userRepo.GetUserByID(fromUserID)
	if err != nil {
		s.log.Error().Err(err).Msgf("Failed to get fromUser with ID: %s", fromUserID)
		return err
	}

	toUser, err := s.userRepo.GetUserByID(toUserID)
	if err != nil {
		s.log.Error().Err(err).Msgf("Failed to get toUser with ID: %s", fromUserID)
		return err
	}

	tx := s.userRepo.BeginTransaction()
	if tx.Error != nil {
		return tx.Error
	}

	if fromUser.Balance < amount {
		tx.Rollback()
		s.log.Info().Msgf("User: %s does not have enough funds", fromUserID)
		return models.ErrInsufficientFunds
	}

	fromUser.Balance -= amount
	toUser.Balance += amount

	transaction := models.Transaction{
		ID:         uuid.New(),
		Type:       models.TransactionTypeTransfer,
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Amount:     amount,
		CreatedAt:  time.Now(),
	}

	if err := s.transactionRepo.CreateTransaction(tx, &transaction); err != nil {
		s.log.Error().Err(err).Msgf("Failed to create transaction")
		return err
	}

	if err := s.userRepo.UpdateUserBalance(tx, fromUserID, fromUser.Balance); err != nil {
		tx.Rollback()
		s.log.Error().Err(err).Msgf("Failed to update balance of user with ID: %s", fromUserID)
		return err
	}

	if err := s.userRepo.UpdateUserBalance(tx, toUserID, toUser.Balance); err != nil {
		tx.Rollback()
		s.log.Error().Err(err).Msgf("Failed to update balance of user with ID: %s", toUserID)
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		s.log.Error().Err(err).Msgf("Failed to commit transaction")
		return err
	}

	s.log.Info().Msgf("Successfully transferred from %s to %s", fromUserID, toUserID)
	return nil
}

func (s *UserServiceImpl) PurchaseMerch(userID uuid.UUID, merchItemName string) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		s.log.Error().Err(err).Msgf("Failed to get user with ID: %s", userID)
		return err
	}

	merchItem, err := s.merchService.GetMerchItemByName(merchItemName)
	if err != nil {
		s.log.Error().Err(err).Msgf("Failed to get merch with name: %s", merchItemName)
		return err
	}

	if user.Balance < merchItem.Price {
		s.log.Info().Msgf("User: %s does not have enough funds", userID)
		return models.ErrNotEnoughFunds
	}

	user.Balance -= merchItem.Price

	transaction := models.Transaction{
		ID:            uuid.New(),
		Type:          models.TransactionTypePurchase,
		FromUserID:    userID,
		ToUserID:      userID,
		Amount:        merchItem.Price,
		MerchItemName: &merchItemName,
		CreatedAt:     time.Now(),
	}

	if err := s.transactionRepo.CreateTransaction(nil, &transaction); err != nil {
		s.log.Error().Err(err).Msgf("Failed to create transaction")
		return err
	}

	if err := s.userRepo.UpdateUserBalance(nil, userID, user.Balance); err != nil {
		s.log.Error().Err(err).Msgf("Failed to update balance of user with ID: %s", userID)
		return err
	}

	s.log.Info().Msgf("Successfully purchased %s by %s", merchItemName, userID)
	return nil
}

func (s *UserServiceImpl) GetPurchasedItems(userID uuid.UUID) ([]models.MerchItem, error) {
	return s.userRepo.GetPurchasedItems(userID)
}

func (s *UserServiceImpl) GetTransactionHistory(userID uuid.UUID) ([]models.Transaction, error) {
	return s.userRepo.GetTransactionHistory(userID)
}

func (s *UserServiceImpl) CreateUser(username, passwordHash string) (*models.User, error) {
	user := &models.User{
		ID:       uuid.New(),
		Username: username,
		Password: passwordHash,
		Balance:  1000,
	}

	err := s.userRepo.CreateUser(user)
	if err != nil {
		s.log.Error().Err(err).Msg("Error creating user")
		return nil, err
	}
	s.log.Info().Msgf("Successfully created user: %s", user.Username)
	return user, nil
}

func (s *UserServiceImpl) GetUserByID(userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			s.log.Info().Err(err).Msg("User not found")
			return nil, err
		}
		s.log.Error().Err(err).Msg("Error getting user by id")
		return nil, err
	}
	return user, nil
}

func (s *UserServiceImpl) GetUserByUsername(username string) (*models.User, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			s.log.Info().Err(err).Msg("User not found")
			return nil, err
		}
		s.log.Error().Err(err).Msg("Error getting user by username")
		return nil, err
	}
	return user, nil
}

func (s *UserServiceImpl) UpdateUser(user *models.User) error {
	err := s.userRepo.UpdateUser(user.ID, user)
	if err != nil {
		s.log.Error().Err(err).Msg("Error updating user")
		return err
	}
	s.log.Info().Msgf("Successfully updated user: %s", user.Username)
	return nil
}

func (s *UserServiceImpl) DeleteUser(userID uuid.UUID) error {
	err := s.userRepo.DeleteUser(userID)
	if err != nil {
		s.log.Error().Err(err).Msg("Error deleting user")
		return err
	}
	s.log.Info().Msgf("Successfully deleted user: %s", userID)
	return nil
}
