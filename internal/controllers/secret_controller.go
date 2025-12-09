package controllers

import (
	"github.com/akansha204/cryptex-secretservice/internal/services"
	"github.com/gofiber/fiber/v2"
)

type SecretController struct {
	service *services.SecretService
}

func NewSecretController(service *services.SecretService) *SecretController {
	return &SecretController{service: service}
}

type CreateSecretBody struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	TTL   *int   `json:"ttl"` // number of days and optional
}

type UpdateSecretBody struct {
	Value *string `json:"value"`
	TTL   *int    `json:"ttl"`
}

func (sc *SecretController) CreateSecret(c *fiber.Ctx) error {

	userID := c.Locals("userId").(string)
	projectID := c.Params("projectId")

	var body CreateSecretBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request json"})
	}

	if body.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name is required"})
	}
	if body.Value == "" {
		return c.Status(400).JSON(fiber.Map{"error": "value is required"})
	}

	secret, err := sc.service.CreateSecret(
		c.Context(),
		userID,
		projectID,
		body.Name,
		body.Value,
		body.TTL,
	)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(secret)
}

func (sc *SecretController) GetSecret(c *fiber.Ctx) error {

	userID := c.Locals("userId").(string)
	projectID := c.Params("projectId")
	secretID := c.Params("secretId")

	secret, plaintext, err := sc.service.GetSecretByID(
		c.Context(),
		userID,
		projectID,
		secretID,
	)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"secret":    secret,
		"plaintext": plaintext,
	})
}

func (sc *SecretController) UpdateSecret(c *fiber.Ctx) error {

	userID := c.Locals("userId").(string)
	projectID := c.Params("projectId")
	secretID := c.Params("secretId")

	var body UpdateSecretBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request json"})
	}

	if body.Value == nil && body.TTL == nil {
		return c.Status(400).JSON(fiber.Map{"error": "nothing to update"})
	}

	updated, err := sc.service.UpdateSecret(
		c.Context(),
		userID,
		projectID,
		secretID,
		body.Value,
		body.TTL,
	)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(updated)
}

func (sc *SecretController) DeleteSecret(c *fiber.Ctx) error {

	userID := c.Locals("userId").(string)
	projectID := c.Params("projectId")
	secretID := c.Params("secretId")

	err := sc.service.DeleteSecret(
		c.Context(),
		userID,
		projectID,
		secretID,
	)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "secret deleted"})
}

func (sc *SecretController) RevokeSecret(c *fiber.Ctx) error {

	userID := c.Locals("userId").(string)
	projectID := c.Params("projectId")
	secretID := c.Params("secretId")

	err := sc.service.RevokeSecret(
		c.Context(),
		userID,
		projectID,
		secretID,
	)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "secret revoked successfully"})
}
