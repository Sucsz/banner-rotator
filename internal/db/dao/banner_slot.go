package dao

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// BannerSlotDAO — интерфейс для работы с таблицей banner_slots (many-to-many).
type BannerSlotDAO interface {
	AddBannerToSlot(ctx context.Context, bannerID, slotID int64) error
	RemoveBannerFromSlot(ctx context.Context, bannerID, slotID int64) error
	GetBannersBySlot(ctx context.Context, slotID int64) ([]int64, error)
	IsBannerInSlot(ctx context.Context, bannerID, slotID int64) (bool, error)
}

type bannerSlotDAO struct {
	conn *pgx.Conn
}

// NewBannerSlotDAO создаёт экземпляр bannerSlotDAO в виде интерфейса BannerSlotDAO
func NewBannerSlotDAO(conn *pgx.Conn) BannerSlotDAO {
	return &bannerSlotDAO{conn: conn}
}

// AddBannerToSlot связывает баннер и слот
func (d *bannerSlotDAO) AddBannerToSlot(ctx context.Context, bannerID, slotID int64) error {
	_, err := d.conn.Exec(ctx, `
        INSERT INTO banner_slots (banner_id, slot_id, created_at)
        VALUES ($1, $2, NOW())
    `, bannerID, slotID)
	if err != nil {
		return fmt.Errorf("BannerSlotDAO.AddBannerToSlot: %w", err)
	}
	return nil
}

// RemoveBannerFromSlot удаляет связь баннера и слота
func (d *bannerSlotDAO) RemoveBannerFromSlot(ctx context.Context, bannerID, slotID int64) error {
	cmd, err := d.conn.Exec(ctx, `
        DELETE FROM banner_slots
        WHERE banner_id = $1 AND slot_id = $2
    `, bannerID, slotID)
	if err != nil {
		return fmt.Errorf("BannerSlotDAO.RemoveBannerFromSlot: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("BannerSlotDAO.RemoveBannerFromSlot: relation (%d,%d) not found", bannerID, slotID)
	}
	return nil
}

// GetBannersBySlot возвращает список banner_id для заданного slot_id
func (d *bannerSlotDAO) GetBannersBySlot(ctx context.Context, slotID int64) ([]int64, error) {
	rows, err := d.conn.Query(ctx, `
        SELECT banner_id
        FROM banner_slots
        WHERE slot_id = $1
        ORDER BY created_at
    `, slotID)
	if err != nil {
		return nil, fmt.Errorf("BannerSlotDAO.GetBannersBySlot: %w", err)
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var bid int64
		if err := rows.Scan(&bid); err != nil {
			return nil, fmt.Errorf("BannerSlotDAO.GetBannersBySlot scan: %w", err)
		}
		ids = append(ids, bid)
	}
	return ids, nil
}

// IsBannerInSlot проверяет, связаны ли баннер и слот
func (d *bannerSlotDAO) IsBannerInSlot(ctx context.Context, bannerID, slotID int64) (bool, error) {
	var exists bool
	err := d.conn.QueryRow(ctx, `
        SELECT EXISTS(
            SELECT 1 FROM banner_slots
            WHERE banner_id = $1 AND slot_id = $2
        )
    `, bannerID, slotID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("BannerSlotDAO.IsBannerInSlot: %w", err)
	}
	return exists, nil
}
