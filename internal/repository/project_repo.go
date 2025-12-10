package repository

import (
	"context"
	"database/sql"
	"time"

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
		Where("deleted_at IS NULL").
		Scan(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}
func (pr *ProjectRepository) GetProjectsByUserID(ctx context.Context, userID string) ([]models.Project, error) {
	var projects []models.Project

	err := database.DB.NewSelect().
		Model(&projects).
		Where("user_id = ?", userID).
		Where("deleted_at IS NULL").
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
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}
func (pr *ProjectRepository) SoftDeleteProject(ctx context.Context, projectID string) error {
	_, err := database.DB.NewUpdate().
		Model(&models.Project{}).
		Set("deleted_at = ?", time.Now()).
		Where("project_id = ?", projectID).
		Exec(ctx)
	return err
}
