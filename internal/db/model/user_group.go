package model

import "time"

// UserGroup — соц-демографическая группа пользователей.
type UserGroup struct {
	ID          int64      `db:"id"`
	Description string     `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at,omitempty"`
}
