package validators

import (
	"inventory-backend/models"
	"regexp"

	"github.com/gofiber/fiber/v2"
)

func ValidateRegister(c *fiber.Ctx) error {
	var user models.RegisterRequest

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate username
	if len(user.Username) < 2 || len(user.Username) > 50 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Username harus antara 2-50 karakter",
		})
	}

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(user.Email) {
		return c.Status(400).JSON(fiber.Map{
			"error": "Format email tidak valid",
		})
	}

	// Validate password
	if len(user.Password) < 6 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Password minimal 6 karakter",
		})
	}

	// Validate role
	if user.Role != "admin" && user.Role != "user" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Role hanya boleh 'admin' atau 'user'",
		})
	}

	// Store validated data in context
	c.Locals("userData", user)
	return c.Next()
}

func ValidateLogin(c *fiber.Ctx) error {
	var loginReq models.LoginRequest

	if err := c.BodyParser(&loginReq); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(loginReq.Email) {
		return c.Status(400).JSON(fiber.Map{
			"error": "Format email tidak valid",
		})
	}

	// Validate password not empty
	if loginReq.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Password tidak boleh kosong",
		})
	}

	// Store validated data in context
	c.Locals("loginData", loginReq)
	return c.Next()
}
