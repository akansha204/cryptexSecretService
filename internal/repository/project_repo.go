package repository

import (
	"context"

	"github.com/akansha204/cryptex-secretservice/internal/database"
	"github.com/akansha204/cryptex-secretservice/internal/models"
)

type ProjectRepository struct{}

func NewProjectRepository() *ProjectRepository {
	return &ProjectRepository{}
}

func (pr *ProjectRepository) CreateProject(ctx context.Context, project *models.Project) error {
	_, err := database.DB.NewInsert().
		Model(project).
		Exec(ctx)
	return err
}

func (pr *ProjectRepository) GetProjectByID(ctx context.Context, projectID string) (*models.Project, error) {
	var project models.Project
	err := database.DB.NewSelect().
		Model(&project).
		Where("project_id = ?", projectID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &project, nil
}
func (pr *ProjectRepository) GetProjectsByUserID(ctx context.Context, userID string) ([]models.Project, error) {
	var projects []models.Project

	err := database.DB.NewSelect().
		Model(&projects).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return projects, nil
}
func (pr *ProjectRepository) UpdateProject(ctx context.Context, project *models.Project) error {
	_, err := database.DB.NewUpdate().
		Model(project).
		Where("project_id = ?", project.ID).
		Exec(ctx)
	return err
}
func (pr *ProjectRepository) DeleteProject(ctx context.Context, projectID string) error {
	_, err := database.DB.NewDelete().
		Model(&models.Project{}).
		Where("project_id = ?", projectID).
		Exec(ctx)
	return err
}
