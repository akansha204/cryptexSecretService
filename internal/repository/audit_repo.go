package repository

import (
	"context"

	"github.com/akansha204/cryptex-secretservice/internal/database"
	"github.com/akansha204/cryptex-secretservice/internal/models"
)

type AuditRepository struct{}

func NewAuditRepository() *AuditRepository {
	return &AuditRepository{}
}

func (r *AuditRepository) Create(ctx context.Context, logEntry *models.AuditLog) error {
	_, err := database.DB.NewInsert().
		Model(logEntry).
		Exec(ctx)
	return err
}

func (r *AuditRepository) FindAll(ctx context.Context) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := database.DB.NewSelect().
		Model(&logs).
		Order("timestamp DESC").
		Scan(ctx)
	return logs, err
}

func (r *AuditRepository) FindByProject(ctx context.Context, projectId string) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := database.DB.NewSelect().
		Model(&logs).
		Where("project_id = ?", projectId).
		Order("timestamp DESC").
		Scan(ctx)
	return logs, err
}

func (r *AuditRepository) FindBySecret(ctx context.Context, secretId string) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := database.DB.NewSelect().
		Model(&logs).
		Where("secret_id = ?", secretId).
		Order("timestamp DESC").
		Scan(ctx)
	return logs, err
}
