package services

import (
	"context"
	"errors"

	"github.com/akansha204/cryptex-secretservice/internal/models"
	"github.com/akansha204/cryptex-secretservice/internal/repository"
	"github.com/google/uuid"
)

type ProjectService struct {
	repo *repository.ProjectRepository
}

func NewProjectService(repo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) CreateProject(ctx context.Context, userID string, name string, description *string) (*models.Project, error) {

	project := &models.Project{
		UserID:      uuid.MustParse(userID),
		Name:        name,
		Description: description,
	}

	err := s.repo.CreateProject(ctx, project)
	if err != nil {
		return nil, err
	}

	return project, nil
}
func (s *ProjectService) GetProjectByID(ctx context.Context, projectID string) (*models.Project, error) {
	return s.repo.GetProjectByID(ctx, projectID)

}
func (s *ProjectService) GetProjectsByUser(ctx context.Context, userID string) ([]models.Project, error) {
	return s.repo.GetProjectsByUserID(ctx, userID)
}

func (s *ProjectService) UpdateProject(ctx context.Context, projectID string, userID string, name string, description *string) (*models.Project, error) {

	project, err := s.repo.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("project not found")
	}
	if project.UserID.String() != userID {
		return nil, errors.New("forbidden: not your project")
	}

	project.Name = name
	project.Description = description

	err = s.repo.UpdateProject(ctx, project)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) DeleteProject(ctx context.Context, projectID string, userID string) error {

	project, err := s.repo.GetProjectByID(ctx, projectID)
	if err != nil {
		return err
	}
	if project == nil {
		return errors.New("project not found")
	}
	if project.UserID.String() != userID {
		return errors.New("forbidden: not your project")
	}

	return s.repo.DeleteProject(ctx, projectID)
}
