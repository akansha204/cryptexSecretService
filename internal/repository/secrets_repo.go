package repository

import (
	"context"
	"database/sql"

	"github.com/akansha204/cryptex-secretservice/internal/database"
	"github.com/akansha204/cryptex-secretservice/internal/models"
)

type SecretRepository struct{}

func NewSecretRepository() *SecretRepository {
	return &SecretRepository{}
}

func (sr *SecretRepository) CreateSecret(ctx context.Context, secret *models.Secret) error {
	_, err := database.DB.NewInsert().
		Model(secret).
		Exec(ctx)
	return err
}
func (sr *SecretRepository) GetSecretByID(ctx context.Context, secretID string) (*models.Secret, error) {
	var secret models.Secret
	err := database.DB.NewSelect().
		Model(&secret).
		Column("*").
		Where("secret_id = ?", secretID).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &secret, nil
}
func (sr *SecretRepository) UpdateSecret(ctx context.Context, secret *models.Secret) error {
	_, err := database.DB.NewUpdate().
		Model(secret).
		Column("s_name", "s_value", "secret_version", "updated_at", "revoked", "ttl", "expires_at").
		Where("secret_id = ?", secret.ID).
		Exec(ctx)
	return err
}
func (sr *SecretRepository) DeleteSecret(ctx context.Context, secretID string) error {
	_, err := database.DB.NewDelete().
		Model(&models.Secret{}).
		Where("secret_id = ?", secretID).
		Exec(ctx)
	return err
}
func (sr *SecretRepository) GetLatestVersion(ctx context.Context, projectID, name string) (*models.Secret, error) {
	var secret models.Secret

	err := database.DB.NewSelect().
		Model(&secret).
		Column("*").
		Where("project_id = ?", projectID).
		Where("s_name = ?", name).
		Order("secret_version DESC").
		Limit(1).
		Scan(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &secret, nil
}
