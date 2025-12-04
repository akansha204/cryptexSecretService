package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type AuditLog struct {
	bun.BaseModel `bun:"table:audit_logs"`

	ID        uuid.UUID `bun:"audit_id,pk,type:uuid,default:gen_random_uuid()"`
	UserID    uuid.UUID `bun:"user_id,type:uuid,nullzero"` // can be null for system events
	ProjectID uuid.UUID `bun:"project_id,type:uuid,nullzero"`
	SecretID  uuid.UUID `bun:"secret_id,type:uuid,nullzero"`

	Action  string  `bun:"action,notnull"`
	Message *string `bun:"message,nullzero"`

	Timestamp time.Time `bun:"timestamp,default:current_timestamp"`
}
