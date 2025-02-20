package unit

import (
	"Avito-backend-trainee-assignment-winter-2025/internal/service"
	"errors"
	"github.com/google/uuid"
	"testing"

	"Avito-backend-trainee-assignment-winter-2025/internal/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMerchRepository struct {
	mock.Mock
}

func (m *MockMerchRepository) GetAll() ([]models.MerchItem, error) {
	args := m.Called()
	return args.Get(0).([]models.MerchItem), args.Error(1)
}

func (m *MockMerchRepository) GetByName(name string) (*models.MerchItem, error) {
	args := m.Called(name)
	return args.Get(0).(*models.MerchItem), args.Error(1)
}

func (m *MockMerchRepository) CreateMerch(item *models.MerchItem) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockMerchRepository) UpdateMerch(name string, item *models.MerchItem) error {
	args := m.Called(name, item)
	return args.Error(0)
}

func (m *MockMerchRepository) DeleteMerch(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *MockMerchRepository) InitMerchItems() error {
	args := m.Called()
	return args.Error(0)
}

func TestMerchService_GetAllMerchItems_Success(t *testing.T) {
	mockRepo := new(MockMerchRepository)
	mockLogger := zerolog.New(zerolog.NewConsoleWriter())
	merchService := service.NewMerchService(mockRepo, &mockLogger)

	mockRepo.On("GetAll").Return([]models.MerchItem{
		{ID: uuid.New(), Name: "T-shirt", Price: 20.0},
		{ID: uuid.New(), Name: "Powerbank", Price: 5.0},
	}, nil)

	items, err := merchService.GetAllMerchItems()

	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Equal(t, "T-shirt", items[0].Name)
	assert.Equal(t, "Powerbank", items[1].Name)
	mockRepo.AssertExpectations(t)
}

func TestMerchService_GetAllMerchItems_Error(t *testing.T) {
	mockRepo := new(MockMerchRepository)
	mockLogger := zerolog.New(zerolog.NewConsoleWriter())
	merchService := service.NewMerchService(mockRepo, &mockLogger)

	mockRepo.On("GetAll").Return([]models.MerchItem(nil), errors.New("database error"))

	items, err := merchService.GetAllMerchItems()

	assert.Error(t, err)
	assert.Nil(t, items)
	mockRepo.AssertExpectations(t)
}

func TestMerchService_GetMerchItemByName_Success(t *testing.T) {
	mockRepo := new(MockMerchRepository)
	mockLogger := zerolog.New(zerolog.NewConsoleWriter())
	merchService := service.NewMerchService(mockRepo, &mockLogger)

	mockRepo.On("GetByName", "T-shirt").Return(&models.MerchItem{ID: uuid.New(), Name: "T-shirt", Price: 20.0}, nil)

	item, err := merchService.GetMerchItemByName("T-shirt")

	assert.NoError(t, err)
	assert.Equal(t, "T-shirt", item.Name)
	mockRepo.AssertExpectations(t)
}

func TestMerchService_GetMerchItemByName_Error(t *testing.T) {
	mockRepo := new(MockMerchRepository)
	mockLogger := zerolog.New(zerolog.NewConsoleWriter())
	merchService := service.NewMerchService(mockRepo, &mockLogger)

	mockRepo.On("GetByName", "T-shirt").Return((*models.MerchItem)(nil), errors.New("item not found"))

	item, err := merchService.GetMerchItemByName("T-shirt")

	assert.Error(t, err)
	assert.Nil(t, item)
	mockRepo.AssertExpectations(t)
}

func TestMerchService_AddMerchItem_Success(t *testing.T) {
	mockRepo := new(MockMerchRepository)
	mockLogger := zerolog.New(zerolog.NewConsoleWriter())
	merchService := service.NewMerchService(mockRepo, &mockLogger)

	item := models.MerchItem{ID: uuid.New(), Name: "T-shirt", Price: 20.0}

	mockRepo.On("CreateMerch", &item).Return(nil)

	err := merchService.AddMerchItem(item)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMerchService_AddMerchItem_Error(t *testing.T) {
	mockRepo := new(MockMerchRepository)
	mockLogger := zerolog.New(zerolog.NewConsoleWriter())
	merchService := service.NewMerchService(mockRepo, &mockLogger)

	item := models.MerchItem{ID: uuid.New(), Name: "T-shirt", Price: 20.0}

	mockRepo.On("CreateMerch", &item).Return(errors.New("failed to add item"))

	err := merchService.AddMerchItem(item)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMerchService_UpdateMerchItem_Success(t *testing.T) {
	mockRepo := new(MockMerchRepository)
	mockLogger := zerolog.New(zerolog.NewConsoleWriter())
	merchService := service.NewMerchService(mockRepo, &mockLogger)

	item := models.MerchItem{ID: uuid.New(), Name: "T-shirt", Price: 25.0}

	mockRepo.On("UpdateMerch", "T-shirt", &item).Return(nil)

	err := merchService.UpdateMerchItem("T-shirt", item)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMerchService_DeleteMerchItem_Success(t *testing.T) {
	mockRepo := new(MockMerchRepository)
	mockLogger := zerolog.New(zerolog.NewConsoleWriter())
	merchService := service.NewMerchService(mockRepo, &mockLogger)

	mockRepo.On("DeleteMerch", "T-shirt").Return(nil)

	err := merchService.DeleteMerchItem("T-shirt")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMerchService_InitMerchItems_Success(t *testing.T) {
	mockRepo := new(MockMerchRepository)
	mockLogger := zerolog.New(zerolog.NewConsoleWriter())
	merchService := service.NewMerchService(mockRepo, &mockLogger)

	mockRepo.On("InitMerchItems").Return(nil)

	err := merchService.InitMerchItems()

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMerchService_InitMerchItems_Error(t *testing.T) {
	mockRepo := new(MockMerchRepository)
	mockLogger := zerolog.New(zerolog.NewConsoleWriter())
	merchService := service.NewMerchService(mockRepo, &mockLogger)

	mockRepo.On("InitMerchItems").Return(errors.New("failed to initialize merch items"))

	err := merchService.InitMerchItems()

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
