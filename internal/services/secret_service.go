package services

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/akansha204/cryptex-secretservice/internal/models"
	"github.com/akansha204/cryptex-secretservice/internal/repository"
	"github.com/akansha204/cryptex-secretservice/internal/utils"
	"github.com/google/uuid"
)

type SecretService struct {
	secretRepo   *repository.SecretRepository
	projectRepo  *repository.ProjectRepository
	AuditService *AuditService
}

func NewSecretService(secretRepo *repository.SecretRepository, projectRepo *repository.ProjectRepository, auditService *AuditService) *SecretService {
	return &SecretService{
		secretRepo:   secretRepo,
		projectRepo:  projectRepo,
		AuditService: auditService,
	}
}

func (s *SecretService) CreateSecret(
	ctx context.Context,
	userID string,
	projectID string,
	name string,
	plaintextValue string,
	ttlDays *int,
) (*models.Secret, error) {

	userUUID := uuid.MustParse(userID)
	project, err := s.projectRepo.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.DeletedAt != nil {
		return nil, errors.New("project not available")
	}
	if project == nil || project.UserID.String() != userID {
		return nil, errors.New("unauthorized: project does not belong to this user")
	}

	latest, err := s.secretRepo.GetLatestVersion(ctx, projectID, name)
	if err != nil {
		return nil, err
	}
	newVersion := 1
	if latest != nil {
		newVersion = latest.Version + 1
	}

	encryptedValue, err := utils.Encrypt(plaintextValue)
	if err != nil {
		return nil, err
	}

	var expiresAt *time.Time
	if ttlDays == nil {
		expiresAt = nil

	} else {
		if *ttlDays < 1 {
			return nil, errors.New("ttl must be at least 1 day")
		}

		t := time.Now().Add(time.Duration(*ttlDays) * 24 * time.Hour)
		expiresAt = &t
	}

	secret := &models.Secret{
		ID:        uuid.New(),
		ProjectID: uuid.MustParse(projectID),
		Name:      name,
		Value:     encryptedValue,
		Version:   newVersion,
		TTL:       ttlDays,
		Revoked:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}

	err = s.secretRepo.CreateSecret(ctx, secret)
	if err != nil {
		return nil, err
	}
	s.AuditService.Log(
		ctx,
		&userUUID,
		&secret.ProjectID,
		&secret.ID,
		"CREATE_SECRET",
		"Secret created with version "+strconv.Itoa(newVersion),
	)

	return secret, nil
}

func (s *SecretService) GetSecretByID(
	ctx context.Context,
	userID string,
	projectID string,
	secretID string,
) (*models.Secret, string, error) {

	project, err := s.projectRepo.GetProjectByID(ctx, projectID)
	if err != nil || project == nil || project.DeletedAt != nil {
		return nil, "", errors.New("project not found")
	}
	if project.UserID.String() != userID {
		return nil, "", errors.New("unauthorized")
	}
	secret, err := s.secretRepo.GetSecretByID(ctx, secretID)
	if err != nil {
		return nil, "", err
	}

	if secret == nil || secret.DeletedAt != nil {
		return nil, "", errors.New("secret not found")
	}

	if secret.ExpiresAt != nil && time.Now().After(*secret.ExpiresAt) {
		return nil, "", errors.New("secret has expired")
	}

	if secret.Revoked {
		return nil, "", errors.New("secret is revoked")
	}

	//decrypt secret value
	plaintext, err := utils.Decrypt(secret.Value)
	if err != nil {
		return nil, "", err
	}

	return secret, plaintext, nil
}

func (s *SecretService) UpdateSecret(
	ctx context.Context,
	userID string,
	projectID string,
	secretID string,
	newValue *string,
	ttlDays *int,
) (*models.Secret, error) {
	userUUID := uuid.MustParse(userID)

	project, err := s.projectRepo.GetProjectByID(ctx, projectID)
	if err != nil || project == nil || project.UserID.String() != userID {
		return nil, errors.New("unauthorized")
	}

	existing, err := s.secretRepo.GetSecretByID(ctx, secretID)
	if err != nil || existing == nil || existing.DeletedAt != nil {
		return nil, errors.New("secret not found")
	}
	if existing.Revoked {
		return nil, errors.New("cannot update a revoked secret")
	}

	valueChanged := false

	if newValue != nil {
		encrypted, err := utils.Encrypt(*newValue)
		if err != nil {
			return nil, err
		}
		existing.Value = encrypted
		valueChanged = true
	}
	if ttlDays != nil {
		if *ttlDays < 1 {
			return nil, errors.New("ttl must be at least 1 day")
		}

		existing.TTL = ttlDays
		expires := time.Now().Add(time.Duration(*ttlDays) * 24 * time.Hour)
		existing.ExpiresAt = &expires
	}

	if valueChanged {
		existing.Version += 1
	}
	existing.UpdatedAt = time.Now()

	err = s.secretRepo.UpdateSecret(ctx, existing)
	if err != nil {
		return nil, err
	}

	s.AuditService.Log(
		ctx,
		&userUUID,
		&existing.ProjectID,
		&existing.ID,
		"UPDATE_SECRET",
		"Secret updated (version "+strconv.Itoa(existing.Version)+")",
	)

	return existing, nil
}

func (s *SecretService) DeleteSecret(
	ctx context.Context,
	userID string,
	projectID string,
	secretID string,
) error {

	userUUID := uuid.MustParse(userID)
	project, err := s.projectRepo.GetProjectByID(ctx, projectID)
	if err != nil || project == nil || project.DeletedAt != nil {
		return errors.New("unauthorized")
	}
	if project.UserID.String() != userID {
		return errors.New("unauthorized")
	}

	secret, err := s.secretRepo.GetSecretByID(ctx, secretID)
	if err != nil || secret == nil || secret.DeletedAt != nil {
		return errors.New("secret not found")
	}
	s.AuditService.Log(
		ctx,
		&userUUID,
		&secret.ProjectID,
		&secret.ID,
		"DELETE_SECRET",
		"Secret deleted",
	)

	return s.secretRepo.SoftDeleteSecret(ctx, secretID)

}

func (s *SecretService) RevokeSecret(
	ctx context.Context,
	userID string,
	projectID string,
	secretID string,
) error {

	userUUID := uuid.MustParse(userID)

	project, err := s.projectRepo.GetProjectByID(ctx, projectID)
	if err != nil || project == nil || project.UserID.String() != userID {
		return errors.New("unauthorized")
	}

	secret, err := s.secretRepo.GetSecretByID(ctx, secretID)
	if err != nil || secret == nil {
		return errors.New("secret not found")
	}
	if secret.DeletedAt != nil {
		return errors.New("secret not found")
	}

	secret.Revoked = true
	secret.UpdatedAt = time.Now()

	err = s.secretRepo.UpdateSecret(ctx, secret)
	if err != nil {
		return err
	}

	s.AuditService.Log(
		ctx,
		&userUUID,
		&secret.ProjectID,
		&secret.ID,
		"REVOKE_SECRET",
		"Secret revoked",
	)

	return nil
}
