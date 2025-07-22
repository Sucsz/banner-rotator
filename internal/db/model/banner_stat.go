package model

import "time"

// BannerStat — агрегированная статистика показов и кликов по группе.
type BannerStat struct {
	BannerID    int64     `db:"banner_id"`
	SlotID      int64     `db:"slot_id"`
	UserGroupID int64     `db:"user_group_id"`
	Impressions int64     `db:"impressions"`
	Clicks      int64     `db:"clicks"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
