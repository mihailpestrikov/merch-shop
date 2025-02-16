package tests

import (
	"Avito-backend-trainee-assignment-winter-2025/internal/models"
	"Avito-backend-trainee-assignment-winter-2025/internal/service"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByID(userID uuid.UUID) (*models.User, error) {
	args := m.Called(userID)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) GetUserByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) UpdateUser(userID uuid.UUID, user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUser(userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUserBalance(tx *gorm.DB, username string, newBalance int) error {
	args := m.Called(tx, username, newBalance)
	return args.Error(0)
}

func (m *MockUserRepository) GetPurchasedItems(username string) ([]models.MerchItem, error) {
	args := m.Called(username)
	return args.Get(0).([]models.MerchItem), args.Error(1)
}

func (m *MockUserRepository) GetTransactionHistory(username string) ([]models.Transaction, error) {
	args := m.Called(username)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) BeginGormTransaction() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockTransactionRepository) CreateTransaction(tx *gorm.DB, transaction *models.Transaction) error {
	args := m.Called(tx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) BeginTransaction() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

type MockMerchService struct {
	mock.Mock
}

func (m *MockMerchService) GetAllMerchItems() ([]models.MerchItem, error) {
	args := m.Called()
	return args.Get(0).([]models.MerchItem), args.Error(1)
}

func (m *MockMerchService) GetMerchItemByName(merchItemName string) (*models.MerchItem, error) {
	args := m.Called(merchItemName)
	if item, ok := args.Get(0).(*models.MerchItem); ok {
		return item, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMerchService) AddMerchItem(item models.MerchItem) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockMerchService) UpdateMerchItem(merchItemName string, item models.MerchItem) error {
	args := m.Called(merchItemName, item)
	return args.Error(0)
}

func (m *MockMerchService) DeleteMerchItem(merchItemName string) error {
	args := m.Called(merchItemName)
	return args.Error(0)
}

func (m *MockMerchService) InitMerchItems() error {
	args := m.Called()
	return args.Error(0)
}

func TestUserService_CreateUser_Success(t *testing.T) {

	mockUserRepo := new(MockUserRepository)
	mockTransactionRepo := new(MockTransactionRepository)
	mockMerchService := new(MockMerchService)
	mockLogger := zerolog.New(zerolog.NewConsoleWriter())
	userService := service.NewUserService(mockUserRepo, mockTransactionRepo, mockMerchService, &mockLogger)

	username := "testuser"
	password := "hashedpassword"

	mockUserRepo.On("CreateUser", mock.Anything).Return(nil)

	user, err := userService.CreateUser(username, password)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, username, user.Username)

	mockUserRepo.AssertExpectations(t)
}

func TestUserService_GetBalance_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTransactionRepo := new(MockTransactionRepository)
	mockMerchService := new(MockMerchService)
	mockLogger := zerolog.New(zerolog.NewConsoleWriter())
	userService := service.NewUserService(mockUserRepo, mockTransactionRepo, mockMerchService, &mockLogger)

	userID := uuid.New()

	mockUserRepo.On("GetUserByID", userID).Return(&models.User{Balance: 1000}, nil)

	balance, err := userService.GetBalance(userID)

	assert.NoError(t, err)
	assert.Equal(t, 1000, balance)

	mockUserRepo.AssertExpectations(t)
}

func TestUserService_PurchaseMerch_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTransactionRepo := new(MockTransactionRepository)
	mockMerchService := new(MockMerchService)
	mockLogger := zerolog.New(zerolog.NewConsoleWriter())
	userService := service.NewUserService(mockUserRepo, mockTransactionRepo, mockMerchService, &mockLogger)

	username := "testuser"
	merchItemName := "T-shirt"

	mockUserRepo.On("GetUserByUsername", username).Return(&models.User{Balance: 500}, nil)

	mockMerchService.On("GetMerchItemByName", merchItemName).Return(&models.MerchItem{Name: merchItemName, Price: 200}, nil)

	mockTransactionRepo.On("CreateTransaction", mock.Anything, mock.Anything).Return(nil)

	mockUserRepo.On("UpdateUserBalance", mock.Anything, username, 300).Return(nil)

	err := userService.PurchaseMerch(username, merchItemName)

	assert.NoError(t, err)

	mockUserRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
	mockMerchService.AssertExpectations(t)
}

func TestUserService_GetInfo_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTransactionRepo := new(MockTransactionRepository)
	mockMerchService := new(MockMerchService)
	mockLogger := zerolog.New(zerolog.NewConsoleWriter())
	userService := service.NewUserService(mockUserRepo, mockTransactionRepo, mockMerchService, &mockLogger)

	userID := uuid.New()
	username := "testuser"

	mockUserRepo.On("GetUserByID", userID).Return(&models.User{Balance: 1000}, nil)

	mockUserRepo.On("GetPurchasedItems", username).Return([]models.MerchItem{{Name: "T-shirt"}}, nil)

	mockUserRepo.On("GetTransactionHistory", username).Return([]models.Transaction{
		{Type: models.TransactionTypeTransfer, FromUsername: "testuser", ToUsername: "anotheruser", Amount: 50},
	}, nil)

	info, err := userService.GetInfo(userID, username)

	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, 1000, info.Coins)
	assert.Len(t, info.Inventory, 1)
	assert.Equal(t, "T-shirt", info.Inventory[0].Type)

	mockUserRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}

func TestUserService_PurchaseMerch_InsufficientBalance(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTransactionRepo := new(MockTransactionRepository)
	mockMerchService := new(MockMerchService)
	mockLogger := zerolog.New(zerolog.NewConsoleWriter())
	userService := service.NewUserService(mockUserRepo, mockTransactionRepo, mockMerchService, &mockLogger)

	username := "user1"
	merchItemName := "T-shirt"
	merchPrice := 1000

	mockUserRepo.On("GetUserByUsername", username).Return(&models.User{Balance: 500}, nil)

	mockMerchService.On("GetMerchItemByName", merchItemName).Return(&models.MerchItem{Price: merchPrice}, nil)

	err := userService.PurchaseMerch(username, merchItemName)

	assert.Equal(t, models.ErrNotEnoughCoins, err)

	mockUserRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
	mockMerchService.AssertExpectations(t)
}
