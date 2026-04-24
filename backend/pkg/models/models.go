package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// User represents a platform user with login credentials.
type User struct {
	bun.BaseModel `bun:"table:users"`

	ID        uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	Email     string    `bun:",unique,not null"`
	Password  string    `bun:",not null"`
	Name      string    ``
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`

	Chats []Chat `bun:"rel:has-many,join:id=user_id"`
}

// Chat represents a conversation session between a user and an LLM.
type Chat struct {
	bun.BaseModel `bun:"table:chats"`

	ID        uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	UserID    uuid.UUID `bun:",type:uuid,not null"`
	Title     string    ``
	Model     string    ``
	History   string    `bun:",type:text"` // Simplified storage of JSON messages
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

