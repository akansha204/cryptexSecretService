package routes

import (
	"github.com/akansha204/cryptex-secretservice/internal/controllers"
	"github.com/akansha204/cryptex-secretservice/internal/middlewares"
	"github.com/akansha204/cryptex-secretservice/internal/repository"
	"github.com/akansha204/cryptex-secretservice/internal/services"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	auditRepo := repository.NewAuditRepository()
	auditService := services.NewAuditService(auditRepo)

	projectRepo := repository.NewProjectRepository()
	projectService := services.NewProjectService(projectRepo, auditService)
	projectController := controllers.NewProjectController(projectService)

	secretRepo := repository.NewSecretRepository()
	secretService := services.NewSecretService(secretRepo, projectRepo, auditService)
	secretController := controllers.NewSecretController(secretService)

	api := app.Group("/api")

	api.Post("/projects", middlewares.GatewayAuth(), projectController.CreateProject)
	api.Get("/projects/:id", middlewares.GatewayAuth(), projectController.GetProject)
	api.Get("/projects", middlewares.GatewayAuth(), projectController.GetUserProjects)
	api.Put("/projects/:id", middlewares.GatewayAuth(), projectController.UpdateProject)
	api.Delete("/projects/:id", middlewares.GatewayAuth(), projectController.DeleteProject)

	secured := api.Group("/projects/:projectId/secrets", middlewares.GatewayAuth())

	secured.Post("/", secretController.CreateSecret)
	secured.Get("/:secretId", secretController.GetSecret)
	secured.Patch("/:secretId", secretController.UpdateSecret)
	secured.Delete("/:secretId", secretController.DeleteSecret)
	secured.Patch("/:secretId/revoke", secretController.RevokeSecret)

}
