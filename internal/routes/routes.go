package routes

import (
	"os"

	"github.com/akansha204/cryptex-secretservice/internal/controllers"
	"github.com/akansha204/cryptex-secretservice/internal/middlewares"
	"github.com/akansha204/cryptex-secretservice/internal/repository"
	"github.com/akansha204/cryptex-secretservice/internal/services"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	jwtSecret := []byte(os.Getenv("jwt.secret"))

	projectRepo := repository.NewProjectRepository()
	projectService := services.NewProjectService(projectRepo)
	projectController := controllers.NewProjectController(projectService)

	secretRepo := repository.NewSecretRepository()
	secretService := services.NewSecretService(secretRepo, projectRepo)
	secretController := controllers.NewSecretController(secretService)

	api := app.Group("/api")

	api.Post("/projects", middlewares.JWTAuth(jwtSecret), projectController.CreateProject)
	api.Get("/projects/:id", middlewares.JWTAuth(jwtSecret), projectController.GetProject)
	api.Get("/projects", middlewares.JWTAuth(jwtSecret), projectController.GetUserProjects)
	api.Put("/projects/:id", middlewares.JWTAuth(jwtSecret), projectController.UpdateProject)
	api.Delete("/projects/:id", middlewares.JWTAuth(jwtSecret), projectController.DeleteProject)

	secured := api.Group("/projects/:projectId/secrets", middlewares.JWTAuth(jwtSecret))

	secured.Post("/", secretController.CreateSecret)
	secured.Get("/:secretId", secretController.GetSecret)
	secured.Patch("/:secretId", secretController.UpdateSecret)
	secured.Delete("/:secretId", secretController.DeleteSecret)
	secured.Patch("/:secretId/revoke", secretController.RevokeSecret)

}
