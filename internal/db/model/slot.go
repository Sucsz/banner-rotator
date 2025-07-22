package model

import "time"

// Slot — место на сайте, где показываем баннеры.
type Slot struct {
	ID          int64      `db:"id"`
	Description string     `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at,omitempty"`
}
