package service

import (
	"Avito-backend-trainee-assignment-winter-2025/internal/models"
	"Avito-backend-trainee-assignment-winter-2025/repository"
	"github.com/google/uuid"
	"time"
)

type UserService interface {
	TransferCoins(fromUserID, toUserID uuid.UUID, amount int) error
	GetBalance(userID uuid.UUID) (int, error)
	PurchaseMerch(userID uuid.UUID, merchItemID uuid.UUID) error
	GetPurchasedItems(userID uuid.UUID) ([]models.MerchItem, error)
	GetTransactionHistory(userID uuid.UUID) ([]models.Transaction, error)
}

type UserServiceImpl struct {
	userRepo     repository.UserRepository
	merchService MerchService
}

func NewUserServiceImpl(userRepo repository.UserRepository, merchService MerchService) *UserServiceImpl {
	return &UserServiceImpl{userRepo: userRepo, merchService: merchService}
}

func (s *UserServiceImpl) GetBalance(userID uuid.UUID) (int, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return 0, err
	}
	return user.Balance, nil
}

func (s *UserServiceImpl) TransferCoins(fromUserID, toUserID uuid.UUID, amount int) error {
	fromUser, err := s.userRepo.GetUserByID(fromUserID)
	if err != nil {
		return err
	}

	toUser, err := s.userRepo.GetUserByID(toUserID)
	if err != nil {
		return err
	}

	if fromUser.Balance < amount {
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

	if err := s.userRepo.UpdateUserBalance(fromUserID, fromUser.Balance, &transaction); err != nil {
		return err
	}

	if err := s.userRepo.UpdateUserBalance(toUserID, toUser.Balance, &transaction); err != nil {
		return err
	}

	return nil
}

func (s *UserServiceImpl) PurchaseMerch(userID uuid.UUID, merchItemID uuid.UUID) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	merchItem, err := s.merchService.GetMerchItemByID(merchItemID)
	if err != nil {
		return err
	}

	if user.Balance < merchItem.Price {
		return models.ErrNotEnoughFunds
	}

	user.Balance -= merchItem.Price

	transaction := models.Transaction{
		ID:          uuid.New(),
		Type:        models.TransactionTypePurchase,
		FromUserID:  userID,
		ToUserID:    userID,
		Amount:      merchItem.Price,
		MerchItemID: &merchItemID,
		CreatedAt:   time.Now(),
	}

	if err := s.userRepo.UpdateUserBalance(userID, user.Balance, &transaction); err != nil {
		return err
	}

	return nil
}

func (s *UserServiceImpl) GetPurchasedItems(userID uuid.UUID) ([]models.MerchItem, error) {
	return s.userRepo.GetPurchasedItems(userID)
}

func (s *UserServiceImpl) GetTransactionHistory(userID uuid.UUID) ([]models.Transaction, error) {
	return s.userRepo.GetTransactionHistory(userID)
}
