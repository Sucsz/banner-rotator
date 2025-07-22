// Package model содержит структуры, соответствующие таблицам базы данных,
// включая Banner, Slot, UserGroup и т.д.
package model

import "time"

// Banner — рекламный баннер.
type Banner struct {
	ID          int64      `db:"id"`
	Title       string     `db:"title"`
	Content     string     `db:"content"`
	Description string     `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at,omitempty"`
}
