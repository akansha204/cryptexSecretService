package middlewares

import "github.com/gofiber/fiber/v2"

func GatewayAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		gatewayHeader := c.Get("X-Gateway-Source")
		if gatewayHeader != "cryptex-gateway" {
			return c.Status(fiber.StatusForbidden).
				JSON(fiber.Map{"error": "direct access forbidden"})
		}

		userId := c.Get("X-User-Id")
		email := c.Get("X-User-Email")

		if userId == "" {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": "missing gateway identity"})
		}

		c.Locals("userId", userId)
		c.Locals("email", email)

		return c.Next()
	}
}
