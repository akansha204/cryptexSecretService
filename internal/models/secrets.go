package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Secret struct {
	bun.BaseModel `bun:"table:secrets"`

	ID        uuid.UUID `bun:"secret_id,pk,type:uuid,default:gen_random_uuid()"`
	ProjectID uuid.UUID `bun:"project_id,type:uuid,notnull"`
	Name      string    `bun:"s_name,notnull"`
	Value     string    `bun:"s_value,notnull"`
	Version   int       `bun:"secret_version,notnull,default:1"`

	TTL       *int       `bun:"ttl,nullzero"` // store minutes/hours as integer
	Revoked   bool       `bun:"revoked,notnull,default:false"`
	CreatedAt time.Time  `bun:"created_at,default:current_timestamp"`
	UpdatedAt time.Time  `bun:"updated_at,default:current_timestamp"`
	ExpiresAt *time.Time `bun:"expires_at,nullzero"` //obtained from TTL and createdAt
	DeletedAt *time.Time `bun:"deleted_at,nullzero"`
}
