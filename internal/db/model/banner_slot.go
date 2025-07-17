package model

import "time"

// BannerSlot — связь баннера и слота (many-to-many)
type BannerSlot struct {
	BannerID  int64     `db:"banner_id"`
	SlotID    int64     `db:"slot_id"`
	CreatedAt time.Time `db:"created_at"`
}
