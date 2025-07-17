package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWT Secret - should match with auth controller
var jwtSecret = []byte("your-secret-key-change-this-in-production")

// JWT Middleware untuk memverifikasi token
func JWTMiddleware(c *fiber.Ctx) error {
	// Get Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{
			"error": "Token tidak ditemukan",
		})
	}

	// Check if header starts with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(401).JSON(fiber.Map{
			"error": "Format token tidak valid",
		})
	}

	// Extract token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse and validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(401, "Metode signing token tidak valid")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Token tidak valid",
		})
	}

	// Check if token is valid
	if !token.Valid {
		return c.Status(401).JSON(fiber.Map{
			"error": "Token tidak valid",
		})
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"error": "Claims token tidak valid",
		})
	}

	// Store user info in context
	c.Locals("user_id", claims["user_id"])
	c.Locals("email", claims["email"])
	c.Locals("role", claims["role"])

	return c.Next()
}

// Middleware untuk role-based access
func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role == nil {
			return c.Status(403).JSON(fiber.Map{
				"error": "Role tidak ditemukan",
			})
		}

		userRole, ok := role.(string)
		if !ok {
			return c.Status(403).JSON(fiber.Map{
				"error": "Role tidak valid",
			})
		}

		// Admin can access everything
		if userRole == "admin" {
			return c.Next()
		}

		// Check if user has required role
		if userRole != requiredRole {
			return c.Status(403).JSON(fiber.Map{
				"error": "Akses ditolak. Role tidak mencukupi",
			})
		}

		return c.Next()
	}
}

// Middleware untuk admin only access
func RequireAdmin(c *fiber.Ctx) error {
	role := c.Locals("role")
	if role == nil {
		return c.Status(403).JSON(fiber.Map{
			"error": "Role tidak ditemukan",
		})
	}

	userRole, ok := role.(string)
	if !ok || userRole != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "Akses ditolak. Hanya admin yang dapat mengakses",
		})
	}

	return c.Next()
}
