package services

import (
	"context"
	"errors"

	"github.com/akansha204/cryptex-secretservice/internal/models"
	"github.com/akansha204/cryptex-secretservice/internal/repository"
	"github.com/google/uuid"
)

type ProjectService struct {
	repo         *repository.ProjectRepository
	AuditService *AuditService
}

func NewProjectService(repo *repository.ProjectRepository, auditService *AuditService) *ProjectService {
	return &ProjectService{
		repo:         repo,
		AuditService: auditService,
	}
}

func (s *ProjectService) CreateProject(ctx context.Context, userID string, name string, description *string) (*models.Project, error) {

	userUUID := uuid.MustParse(userID)

	project := &models.Project{
		UserID:      userUUID,
		Name:        name,
		Description: description,
	}

	err := s.repo.CreateProject(ctx, project)
	if err != nil {
		return nil, err
	}

	s.AuditService.Log(
		ctx,
		&userUUID,
		&project.ID,
		nil,
		"CREATE_PROJECT",
		"Project created successfully",
	)

	return project, nil
}
func (s *ProjectService) GetProjectByID(ctx context.Context, projectID string) (*models.Project, error) {
	project, err := s.repo.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project == nil || project.DeletedAt != nil {
		return nil, errors.New("project not found")
	}
	return project, nil

}
func (s *ProjectService) GetProjectsByUser(ctx context.Context, userID string) ([]models.Project, error) {
	return s.repo.GetProjectsByUserID(ctx, userID)
}

func (s *ProjectService) UpdateProject(ctx context.Context, projectID string, userID string, name string, description *string) (*models.Project, error) {
	userUUID := uuid.MustParse(userID)

	project, err := s.repo.GetProjectByID(ctx, projectID)
	if err != nil || project == nil {
		return nil, errors.New("project not found")
	}
	if project.DeletedAt != nil {
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

	s.AuditService.Log(
		ctx,
		&userUUID,
		&project.ID,
		nil,
		"UPDATE_PROJECT",
		"Project details updated",
	)
	return project, nil
}

func (s *ProjectService) DeleteProject(ctx context.Context, projectID string, userID string) error {
	userUUID := uuid.MustParse(userID)

	project, err := s.repo.GetProjectByID(ctx, projectID)
	if err != nil || project == nil {
		return errors.New("project not found")
	}
	if project.DeletedAt != nil {
		return errors.New("project already deleted")
	}

	if project.UserID.String() != userID {
		return errors.New("forbidden: not your project")
	}
	s.AuditService.Log(
		ctx,
		&userUUID,
		&project.ID,
		nil,
		"DELETE_PROJECT",
		"Project deleted successfully",
	)

	return s.repo.SoftDeleteProject(ctx, projectID)
}
