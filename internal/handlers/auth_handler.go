package handlers

import (
	"Avito-backend-trainee-assignment-winter-2025/internal/dto"
	"Avito-backend-trainee-assignment-winter-2025/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type AuthHandler interface {
	Authenticate(c *fiber.Ctx) error
}

type AuthHandlerImpl struct {
	authService service.AuthService
	userService service.UserService
	log         *zerolog.Logger
}

func NewAuthHandler(authService service.AuthService, userService service.UserService, log *zerolog.Logger) *AuthHandlerImpl {
	return &AuthHandlerImpl{authService: authService, userService: userService, log: log}
}

// Authenticate godoc
// @Summary Аутентификация и получение JWT-токена
// @Description Позволяет пользователю аутентифицироваться и получить JWT-токен. Если пользователь не существует, он будет автоматически зарегистрирован.
// @Tags auth
// @Accept json
// @Produce json
// @Param req body dto.AuthRequest true "Тело запроса для авторизации"
// @Success 200 {object} dto.AuthResponse "Успешная аутентификация"
// @Failure 400 {object} dto.ErrorResponse "Неверный запрос"
// @Failure 401 {object} dto.ErrorResponse "Неавторизован"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/auth [post]
func (h *AuthHandlerImpl) Authenticate(c *fiber.Ctx) error {
	var req dto.AuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "Invalid request body",
		})
	}

	token, err := h.authService.LoginOrRegister(req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	h.log.Info().Msg("User logged in")

	return c.Status(fiber.StatusOK).JSON(dto.AuthResponse{
		Token: token,
	})
}
