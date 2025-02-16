package service

import (
	"Avito-backend-trainee-assignment-winter-2025/internal/dto"
	"Avito-backend-trainee-assignment-winter-2025/internal/models"
	"Avito-backend-trainee-assignment-winter-2025/internal/repository"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type UserService interface {
	SendCoins(fromUsername, toUsername string, amount int) error
	GetBalance(userID uuid.UUID) (int, error)
	PurchaseMerch(username string, merchItemName string) error
	GetInfo(userID uuid.UUID, username string) (*dto.InfoResponse, error)
	CreateUser(username, password string) (*models.User, error)
	GetUserByID(userID uuid.UUID) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(userID uuid.UUID) error
	UserExists(username string) (bool, error)
	AddCoins(userID uuid.UUID, amount int) error
}

type UserServiceImpl struct {
	userRepo        repository.UserRepository
	transactionRepo repository.TransactionRepository
	merchService    MerchService
	log             *zerolog.Logger
}

func NewUserService(userRepo repository.UserRepository, transactionRepo repository.TransactionRepository, merchService MerchService, log *zerolog.Logger) *UserServiceImpl {
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

func (s *UserServiceImpl) SendCoins(fromUsername, toUsername string, amount int) error {
	fromUser, err := s.userRepo.GetUserByUsername(fromUsername)
	if err != nil {
		s.log.Error().Err(err).Msgf("Failed to get fromUser with username: %s", fromUsername)

		return err
	}

	toUser, err := s.userRepo.GetUserByUsername(toUsername)
	if err != nil {
		s.log.Error().Err(err).Msgf("Failed to get toUser with username: %s", fromUsername)

		return err
	}

	tx := s.transactionRepo.BeginGormTransaction()
	if tx.Error != nil {
		return tx.Error
	}

	if fromUser.Balance < amount {
		tx.Rollback()
		s.log.Info().Msgf("User: %s does not have enough coins", fromUsername)

		return models.ErrNotEnoughCoins
	}

	fromUser.Balance -= amount
	toUser.Balance += amount

	transaction := models.Transaction{
		ID:           uuid.New(),
		Type:         models.TransactionTypeTransfer,
		FromUsername: fromUsername,
		ToUsername:   toUsername,
		Amount:       amount,
		CreatedAt:    time.Now(),
	}

	if err := s.transactionRepo.CreateTransaction(tx, &transaction); err != nil {
		s.log.Error().Err(err).Msgf("Failed to create transaction")

		return err
	}

	if err := s.userRepo.UpdateUserBalance(tx, fromUsername, fromUser.Balance); err != nil {
		tx.Rollback()
		s.log.Error().Err(err).Msgf("Failed to update balance of user with username: %s", fromUsername)

		return err
	}

	if err := s.userRepo.UpdateUserBalance(tx, toUsername, toUser.Balance); err != nil {
		tx.Rollback()
		s.log.Error().Err(err).Msgf("Failed to update balance of user with username: %s", toUsername)

		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		s.log.Error().Err(err).Msgf("Failed to commit transaction")

		return err
	}

	s.log.Info().Int("amount", amount).Msgf("Successfully transferred from %s to %s", fromUsername, toUsername)

	return nil
}

func (s *UserServiceImpl) PurchaseMerch(username string, merchItemName string) error {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		s.log.Error().Err(err).Msgf("Failed to get user with username: %s", username)

		return err
	}

	merchItem, err := s.merchService.GetMerchItemByName(merchItemName)
	if err != nil {
		s.log.Error().Err(err).Msgf("Failed to get merch with name: %s", merchItemName)

		return err
	}

	if user.Balance < merchItem.Price {
		s.log.Info().Msgf("User: %s does not have enough coins", username)

		return models.ErrNotEnoughCoins
	}

	user.Balance -= merchItem.Price

	transaction := models.Transaction{
		ID:            uuid.New(),
		Type:          models.TransactionTypePurchase,
		FromUsername:  username,
		ToUsername:    username,
		Amount:        merchItem.Price,
		MerchItemName: &merchItemName,
		CreatedAt:     time.Now(),
	}

	if err := s.transactionRepo.CreateTransaction(nil, &transaction); err != nil {
		s.log.Error().Err(err).Msgf("Failed to create transaction")

		return err
	}

	if err := s.userRepo.UpdateUserBalance(nil, username, user.Balance); err != nil {
		s.log.Error().Err(err).Msgf("Failed to update balance of user with username: %s", username)

		return err
	}

	s.log.Info().Msgf("Successfully purchased %s by %s", merchItemName, username)

	return nil
}

func (s *UserServiceImpl) GetInfo(userID uuid.UUID, username string) (*dto.InfoResponse, error) {
	balance, err := s.GetBalance(userID)
	if err != nil {
		s.log.Error().Err(err).Msg("Error getting user balance")

		return nil, err
	}

	purchased, err := s.userRepo.GetPurchasedItems(username)
	if err != nil {
		s.log.Error().Err(err).Msg("Error getting user inventory")

		return nil, err
	}

	itemCounts := make(map[string]int)
	for _, item := range purchased {
		itemCounts[item.Name]++
	}

	var inventory []dto.InventoryItem
	for itemName, quantity := range itemCounts {
		inventory = append(inventory, dto.InventoryItem{
			Type:     itemName,
			Quantity: quantity,
		})
	}

	history, err := s.userRepo.GetTransactionHistory(username)
	if err != nil {
		s.log.Error().Err(err).Msg("Error getting user transaction history")

		return nil, err
	}

	var sent []dto.SentCoin
	var received []dto.ReceivedCoin
	for _, tx := range history {
		if tx.Type == models.TransactionTypeTransfer {
			if tx.FromUsername == username {
				sent = append(sent, dto.SentCoin{
					ToUser: tx.ToUsername,
					Amount: tx.Amount,
				})
			} else if tx.ToUsername == username {
				received = append(received, dto.ReceivedCoin{
					FromUser: tx.FromUsername,
					Amount:   tx.Amount,
				})
			}
		}
	}

	info := &dto.InfoResponse{
		Coins:     balance,
		Inventory: inventory,
		CoinHistory: dto.CoinHistory{
			Sent:     sent,
			Received: received,
		},
	}

	s.log.Info().Msg("Successfully retrieved user inventory")

	return info, nil
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

func (s *UserServiceImpl) UserExists(username string) (bool, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return false, err
	}

	return user != nil, nil
}

func (s *UserServiceImpl) AddCoins(userID uuid.UUID, amount int) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	user.Balance += amount
	err = s.userRepo.UpdateUser(user.ID, user)
	if err != nil {
		return err
	}

	s.log.Info().Int("amount", amount).Msg("Successfully added coins")

	return nil
}
