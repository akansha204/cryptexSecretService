package controllers

import (
	"github.com/akansha204/cryptex-secretservice/internal/services"
	"github.com/gofiber/fiber/v2"
)

type ProjectController struct {
	service *services.ProjectService
}

func NewProjectController(service *services.ProjectService) *ProjectController {
	return &ProjectController{service: service}
}

type Projectbody struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

func (pc *ProjectController) CreateProject(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)

	var body Projectbody

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid json"})
	}
	if body.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name is required"})
	}
	project, err := pc.service.CreateProject(c.Context(), userID, body.Name, body.Description)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(project)
}

func (pc *ProjectController) GetProject(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	projectID := c.Params("id")
	project, err := pc.service.GetProjectByID(c.Context(), projectID)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if project == nil {
		return c.Status(404).JSON(fiber.Map{"error": "project not found"})
	}
	if project.UserID.String() != userID {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden: you do not own this project"})
	}

	return c.JSON(project)
}
func (pc *ProjectController) GetUserProjects(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)

	projects, err := pc.service.GetProjectsByUser(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(projects)
}

func (pc *ProjectController) UpdateProject(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	projectID := c.Params("id")

	var body Projectbody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid json"})
	}

	project, err := pc.service.UpdateProject(c.Context(), projectID, userID, body.Name, body.Description)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(project)
}

func (pc *ProjectController) DeleteProject(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	projectID := c.Params("id")

	err := pc.service.DeleteProject(c.Context(), projectID, userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "project deleted"})
}
