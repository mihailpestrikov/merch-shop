package handlers

import (
	"Avito-backend-trainee-assignment-winter-2025/internal/dto"
	"Avito-backend-trainee-assignment-winter-2025/internal/service"

	"github.com/gofiber/fiber/v2"
)

type MerchHandler struct {
	merchService service.MerchService
}

func NewMerchHandler(merchService service.MerchService) *MerchHandler {
	return &MerchHandler{merchService: merchService}
}

// GetAllMerch godoc
// @Summary Получить список всех товаров
// @Description Возвращает список всех доступных товаров с их названиями и ценами.
// @Tags merch
// @Accept json
// @Produce json
// @Success 200 {object} dto.MerchItemsResponse "Успешный ответ"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/merch [get]
func (h *MerchHandler) GetAllMerch(c *fiber.Ctx) error {
	merchItems, err := h.merchService.GetAllMerchItems()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: "Internal server error",
		})
	}

	var response dto.MerchItemsResponse
	for _, item := range merchItems {
		response.MerchItems = append(response.MerchItems, dto.MerchItemResponse{
			Name:  item.Name,
			Price: item.Price,
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetMerchByName godoc
// @Summary Получить информацию о товаре по имени
// @Description Позволяет получить информацию о товаре по его имени.
// @Tags merch
// @Accept json
// @Produce json
// @Param name path string true "Имя товара"
// @Success 200 {object} dto.MerchItemResponse "Успешный ответ"
// @Failure 400 {object} dto.ErrorResponse "Неверный запрос"
// @Failure 404 {object} dto.ErrorResponse "Товар не найден"
// @Router /api/merch/{name} [get]
func (h *MerchHandler) GetMerchByName(c *fiber.Ctx) error {
	name := c.Params("name")
	if name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "Merch name is required",
		})
	}

	merchItem, err := h.merchService.GetMerchItemByName(name)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	response := dto.MerchItemResponse{
		Name:  merchItem.Name,
		Price: merchItem.Price,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
