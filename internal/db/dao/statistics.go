package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/Sucsz/banner-rotator/internal/db/model"
	"github.com/jackc/pgx/v5"
)

// StatDAO — интерфейс для работы с агрегированной статистикой banner_stats.
type StatDAO interface {
	IncrementView(ctx context.Context, slotID, bannerID, groupID int64) error
	IncrementClick(ctx context.Context, slotID, bannerID, groupID int64) error
	Get(ctx context.Context, slotID, bannerID, groupID int64) (*model.BannerStat, error)
}

type statDAO struct {
	conn *pgx.Conn
}

// NewStatDAO создаёт экземпляр statDAO в виде интерфейса StatDAO.
func NewStatDAO(conn *pgx.Conn) StatDAO {
	return &statDAO{conn: conn}
}

// IncrementView прибавляет 1 к полю impressions, либо создаёт запись.
func (d *statDAO) IncrementView(ctx context.Context, slotID, bannerID, groupID int64) error {
	_, err := d.conn.Exec(ctx, `
        INSERT INTO banner_stats (banner_id, slot_id, user_group_id, impressions, clicks, created_at, updated_at)
        VALUES ($1, $2, $3, 1, 0, NOW(), NOW())
        ON CONFLICT (banner_id, slot_id, user_group_id) DO
          UPDATE SET impressions = banner_stats.impressions + 1,
                     updated_at = NOW()
    `, bannerID, slotID, groupID)
	if err != nil {
		return fmt.Errorf("StatDAO.IncrementView: %w", err)
	}
	return nil
}

// IncrementClick прибавляет 1 к полю clicks, либо создаёт запись.
func (d *statDAO) IncrementClick(ctx context.Context, slotID, bannerID, groupID int64) error {
	_, err := d.conn.Exec(ctx, `
        INSERT INTO banner_stats (banner_id, slot_id, user_group_id, impressions, clicks, created_at, updated_at)
        VALUES ($1, $2, $3, 0, 1, NOW(), NOW())
        ON CONFLICT (banner_id, slot_id, user_group_id) DO
          UPDATE SET clicks = banner_stats.clicks + 1,
                     updated_at = NOW()
    `, bannerID, slotID, groupID)
	if err != nil {
		return fmt.Errorf("StatDAO.IncrementClick: %w", err)
	}
	return nil
}

// Get возвращает агрегированную статистику по тройке ключей.
func (d *statDAO) Get(ctx context.Context, slotID, bannerID, groupID int64) (*model.BannerStat, error) {
	row := d.conn.QueryRow(ctx, `
        SELECT banner_id, slot_id, user_group_id, impressions, clicks, created_at, updated_at
        FROM banner_stats
        WHERE banner_id = $1 AND slot_id = $2 AND user_group_id = $3
    `, bannerID, slotID, groupID)

	var s model.BannerStat
	err := row.Scan(
		&s.BannerID,
		&s.SlotID,
		&s.UserGroupID,
		&s.Impressions,
		&s.Clicks,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("StatDAO.Get: %w", err)
	}
	return &s, nil
}
