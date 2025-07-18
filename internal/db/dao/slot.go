package dao

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Sucsz/banner-rotator/internal/db/model"
	"github.com/jackc/pgx/v5"
)

// SlotDAO — интерфейс для работы с таблицей slots
type SlotDAO interface {
	Create(ctx context.Context, slot *model.Slot) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.Slot, error)
	List(ctx context.Context) ([]model.Slot, error)
	Delete(ctx context.Context, id int64) error
	SoftDelete(ctx context.Context, id int64) error
	Update(ctx context.Context, slot *model.Slot) error
}

type slotDAO struct {
	conn *pgx.Conn
}

// NewSlotDAO создаёт экземпляр slotDAO в виде интерфейса SlotDAO
func NewSlotDAO(conn *pgx.Conn) SlotDAO {
	return &slotDAO{conn: conn}
}

// Create вставляет новую запись и возвращает её ID
func (d *slotDAO) Create(ctx context.Context, slot *model.Slot) (int64, error) {
	var id int64
	now := time.Now()
	err := d.conn.QueryRow(ctx, `
        INSERT INTO slots (description, created_at, updated_at)
        VALUES ($1, $2, $3)
        RETURNING id
    `, slot.Description, now, now).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("SlotDAO.Create: %w", err)
	}
	return id, nil
}

// GetByID возвращает слот по ID, исключая soft-deleted
func (d *slotDAO) GetByID(ctx context.Context, id int64) (*model.Slot, error) {
	row := d.conn.QueryRow(ctx, `
        SELECT id, description, created_at, updated_at, deleted_at
        FROM slots
        WHERE id = $1 AND deleted_at IS NULL
    `, id)

	var s model.Slot
	err := row.Scan(
		&s.ID,
		&s.Description,
		&s.CreatedAt,
		&s.UpdatedAt,
		&s.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("SlotDAO.GetByID: %w", err)
	}
	return &s, nil
}

// List возвращает все не soft-deleted слоты
func (d *slotDAO) List(ctx context.Context) ([]model.Slot, error) {
	rows, err := d.conn.Query(ctx, `
        SELECT id, description, created_at, updated_at, deleted_at
        FROM slots
        WHERE deleted_at IS NULL
        ORDER BY id
    `)
	if err != nil {
		return nil, fmt.Errorf("SlotDAO.List: %w", err)
	}
	defer rows.Close()

	var out []model.Slot
	for rows.Next() {
		var s model.Slot
		if err := rows.Scan(
			&s.ID,
			&s.Description,
			&s.CreatedAt,
			&s.UpdatedAt,
			&s.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("SlotDAO.List scan: %w", err)
		}
		out = append(out, s)
	}
	return out, nil
}

// Delete физически удаляет запись
func (d *slotDAO) Delete(ctx context.Context, id int64) error {
	cmd, err := d.conn.Exec(ctx, `DELETE FROM slots WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("SlotDAO.Delete: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("SlotDAO.Delete: slot %d not found", id)
	}
	return nil
}

// SoftDelete выставляет DeletedAt = now() для метки удаления
func (d *slotDAO) SoftDelete(ctx context.Context, id int64) error {
	cmd, err := d.conn.Exec(ctx, `
        UPDATE slots
        SET deleted_at = $1
        WHERE id = $2 AND deleted_at IS NULL
    `, time.Now(), id)
	if err != nil {
		return fmt.Errorf("SlotDAO.SoftDelete: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("SlotDAO.SoftDelete: slot %d not found or already deleted", id)
	}
	return nil
}

// Update обновляет описание и UpdatedAt
func (d *slotDAO) Update(ctx context.Context, slot *model.Slot) error {
	cmd, err := d.conn.Exec(ctx, `
        UPDATE slots
        SET description = $1, updated_at = $2
        WHERE id = $3 AND deleted_at IS NULL
    `, slot.Description, time.Now(), slot.ID)
	if err != nil {
		return fmt.Errorf("SlotDAO.Update: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("SlotDAO.Update: slot %d not found or deleted", slot.ID)
	}
	return nil
}
