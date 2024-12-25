package middlewares

import (
	"github.com/XzerozZ/Kasian_Phrom_BE/configs"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware(config configs.JWT) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tokenString := ctx.Get("Authorization")
		if tokenString == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":      "Unauthorized",
				"status_code": fiber.StatusUnauthorized,
				"message":     "Missing or invalid token",
				"result":      nil,
			})
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(config.Secret), nil
		})

		if err != nil || !token.Valid {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":      "Unauthorized",
				"status_code": fiber.StatusUnauthorized,
				"message":     "Invalid or expired token",
				"result":      nil,
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":      "Unauthorized",
				"status_code": fiber.StatusUnauthorized,
				"message":     "Invalid token claims",
				"result":      nil,
			})
		}

		ctx.Locals("user_id", claims["user_id"])
		ctx.Locals("role", claims["role"])
		return ctx.Next()
	}
}