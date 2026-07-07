package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/zunarvy19/finance-backend/configs"
	"github.com/zunarvy19/finance-backend/internal/common"
	"github.com/zunarvy19/finance-backend/pkg/jwt"
)

func AuthMiddleware(cfg *configs.Config) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return common.SendError(c, fiber.StatusUnauthorized, "Missing authorization header", nil)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return common.SendError(c, fiber.StatusUnauthorized, "Invalid authorization format", nil)
		}

		tokenString := parts[1]
		claims, err := jwt.ValidateToken(tokenString, cfg.JWTSecret)
		if err != nil {
			return common.SendError(c, fiber.StatusUnauthorized, "Invalid or expired token", nil)
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)

		return c.Next()
	}
}
