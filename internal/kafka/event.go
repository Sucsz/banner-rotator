// Package kafka содержит константы и утилиты для работы с событиями Kafka.
package kafka

import "time"

// EventType — тип события (клик или показ).
type EventType string

const (
	// EventClick — тип события “клик” для Kafka.
	EventClick EventType = "click"
	// EventView — тип события “просмотр” для Kafka.
	EventView EventType = "view"
)

// BannerEvent — структура события для Kafka.
type BannerEvent struct {
	Type        EventType `json:"type"`
	SlotID      int64     `json:"slot_id"`
	BannerID    int64     `json:"banner_id"`
	UserGroupID int64     `json:"user_group_id"`
	Timestamp   time.Time `json:"timestamp"`
}
