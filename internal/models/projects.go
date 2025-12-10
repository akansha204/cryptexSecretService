package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Project struct {
	bun.BaseModel `bun:"table:projects"`

	ID          uuid.UUID `bun:"project_id,pk,type:uuid,default:gen_random_uuid()"`
	UserID      uuid.UUID `bun:"user_id,type:uuid,notnull"`
	Name        string    `bun:"project_name,notnull"`
	Description *string   `bun:"p_description,nullzero"`

	CreatedAt time.Time  `bun:"created_at,default:current_timestamp"`
	UpdatedAt time.Time  `bun:"updated_at,default:current_timestamp"`
	DeletedAt *time.Time `bun:"deleted_at,nullzero"`
}
