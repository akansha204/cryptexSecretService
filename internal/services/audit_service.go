package services

import (
	"context"
	"log"
	"time"

	"github.com/akansha204/cryptex-secretservice/internal/models"
	"github.com/akansha204/cryptex-secretservice/internal/repository"
	"github.com/google/uuid"
)

type AuditService struct {
	repo *repository.AuditRepository
}

func NewAuditService(repo *repository.AuditRepository) *AuditService {
	return &AuditService{
		repo: repo,
	}
}

// ------------------------------------------------------------
// Log an audit event
// ------------------------------------------------------------
func (a *AuditService) Log(
	ctx context.Context,
	userId, projectId, secretId *uuid.UUID,
	action string,
	message string,
) {

	logEntry := &models.AuditLog{
		UserID:    valueOrNil(userId),    // convert nil → uuid.Nil
		ProjectID: valueOrNil(projectId), // convert nil → uuid.Nil
		SecretID:  valueOrNil(secretId),  // convert nil → uuid.Nil
		Action:    action,
		Message:   &message,
		Timestamp: time.Now(),
	}

	if err := a.repo.Create(ctx, logEntry); err != nil {
		log.Printf("[AUDIT ERROR] Failed to insert log: %v", err)
		return
	}

	// Terminal logging (safe — no sensitive values printed)
	log.Printf(
		"[AUDIT] action=%s user=%s project=%s secret=%s message=\"%s\"",
		action,
		logEntry.UserID.String(),
		logEntry.ProjectID.String(),
		logEntry.SecretID.String(),
		message,
	)
}

// ------------------------------------------------------------
// Utility: convert *uuid.UUID to uuid.UUID
// ------------------------------------------------------------
func valueOrNil(id *uuid.UUID) uuid.UUID {
	if id == nil {
		return uuid.Nil
	}
	return *id
}
