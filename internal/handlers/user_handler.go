package handlers

import (
	"Avito-backend-trainee-assignment-winter-2025/internal/dto"
	"Avito-backend-trainee-assignment-winter-2025/internal/models"
	"Avito-backend-trainee-assignment-winter-2025/internal/service"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	authService service.AuthService
	userService service.UserService
}

func NewUserHandler(authService service.AuthService, userService service.UserService) *UserHandler {
	return &UserHandler{authService: authService, userService: userService}
}

// BuyItem godoc
// @Summary Купить предмет за монеты
// @Description Позволяет пользователю купить предмет за монеты.
// @Tags user
// @Accept json
// @Produce json
// @Param item path string true "Название предмета"
// @Success 200 {object} map[string]interface{} "Успешный ответ"
// @Failure 401 {object} dto.ErrorResponse "Неавторизован"
// @Failure 400 {object} dto.ErrorResponse "Недостаточно монет"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/buy/{item} [get]
func (h *UserHandler) BuyItem(c *fiber.Ctx) error {
	itemName := c.Params("item")

	username, ok := c.Locals("username").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error: "Unauthorized",
		})
	}

	err := h.userService.PurchaseMerch(username, itemName)
	if err != nil {
		if errors.Is(err, models.ErrNotEnoughCoins) {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
				Error: "Not enough coins",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: "Internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Item purchased successfully",
	})
}

// SendCoins godoc
// @Summary Отправить монеты другому пользователю
// @Description Позволяет пользователю отправить монеты другому пользователю.
// @Tags user
// @Accept json
// @Produce json
// @Param request body dto.SendCoinRequest true "Тело запроса для передачи монет"
// @Success 200 {string} string "Успешный ответ"
// @Failure 401 {object} dto.ErrorResponse "Неавторизован"
// @Failure 400 {object} dto.ErrorResponse "Неверный запрос"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/send-coins [post]
func (h *UserHandler) SendCoins(c *fiber.Ctx) error {
	username, ok := c.Locals("username").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error: "Unauthorized",
		})
	}

	var request dto.SendCoinRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "Invalid request body",
		})
	}

	err := h.userService.SendCoins(username, request.ToUser, request.Amount)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).SendString("Coins sent successfully")
}

// GetInfo godoc
// @Summary Получить информацию о монетах, инвентаре и истории транзакций
// @Description Позволяет пользователю получить информацию о своих монетах, инвентаре и истории транзакций.
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} dto.InfoResponse "Успешный ответ"
// @Failure 401 {object} dto.ErrorResponse "Неавторизован"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/info [get]
func (h *UserHandler) GetInfo(c *fiber.Ctx) error {
	username, ok1 := c.Locals("username").(string)
	userID, ok2 := c.Locals("userID").(uuid.UUID)

	if !ok1 || !ok2 {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error: "Unauthorized",
		})
	}

	userInfo, err := h.userService.GetInfo(userID, username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: "Failed to get user info",
		})
	}

	return c.Status(fiber.StatusOK).JSON(userInfo)
}

// AddCoins godoc
// @Summary Добавление монет пользователю
// @Description Увеличивает баланс монет у пользователя, если передан корректный `amount`.
// @Tags user
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param amount path int true "Количество монет"
// @Success 200 {string} string "Coins added successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid amount parameter"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Failed to add coins"
// @Router /api/add-coins/{amount} [post]
func (h *UserHandler) AddCoins(c *fiber.Ctx) error {
	amountStr := c.Params("amount")

	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid amount parameter",
		})
	}

	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error: "Unauthorized",
		})
	}

	if err := h.userService.AddCoins(userID, amount); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: "Failed to add coins",
		})
	}

	return c.Status(fiber.StatusOK).SendString("Coins added successfully")
}
