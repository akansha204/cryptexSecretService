package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/akansha204/cryptex-secretservice/internal/database"
)

type PurgeRepository struct{}

func NewPurgeRepository() *PurgeRepository {
	return &PurgeRepository{}
}

// Permanently delete soft-deleted records older than X days
func (r *PurgeRepository) PurgeOldData(ctx context.Context, olderThan time.Duration) error {

	threshold := time.Now().Add(-olderThan)

	_, err := database.DB.NewDelete().
		TableExpr("secrets").
		Where("deleted_at IS NOT NULL").
		Where("deleted_at < ?", threshold).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to purge secrets: %w", err)
	}

	_, err = database.DB.NewDelete().
		TableExpr("projects").
		Where("deleted_at IS NOT NULL").
		Where("deleted_at < ?", threshold).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to purge projects: %w", err)
	}

	return nil
}
