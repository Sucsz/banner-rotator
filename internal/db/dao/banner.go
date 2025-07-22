package dao

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Sucsz/banner-rotator/internal/db/model"
	"github.com/jackc/pgx/v5"
)

// BannerDAO — интерфейс для работы с таблицей banners.
type BannerDAO interface {
	Create(ctx context.Context, banner *model.Banner) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.Banner, error)
	List(ctx context.Context) ([]model.Banner, error)
	Delete(ctx context.Context, id int64) error
	SoftDelete(ctx context.Context, id int64) error
	Update(ctx context.Context, banner *model.Banner) error
}

type bannerDAO struct {
	conn *pgx.Conn
}

// NewBannerDAO создаёт экземпляр bannerDAO в виде интерфейса BannerDAO.
func NewBannerDAO(conn *pgx.Conn) BannerDAO {
	return &bannerDAO{conn: conn}
}

// Create вставляет новую запись и возвращает её ID.
func (d *bannerDAO) Create(ctx context.Context, banner *model.Banner) (int64, error) {
	var id int64
	now := time.Now()
	err := d.conn.QueryRow(ctx, `
        INSERT INTO banners (title, content, description, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `, banner.Title, banner.Content, banner.Description, now, now).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("BannerDAO.Create: %w", err)
	}
	return id, nil
}

// GetByID возвращает баннер по ID, исключая soft-deleted.
func (d *bannerDAO) GetByID(ctx context.Context, id int64) (*model.Banner, error) {
	row := d.conn.QueryRow(ctx, `
        SELECT id, title, content, description, created_at, updated_at, deleted_at
        FROM banners
        WHERE id = $1 AND deleted_at IS NULL
    `, id)

	var b model.Banner
	err := row.Scan(
		&b.ID,
		&b.Title,
		&b.Content,
		&b.Description,
		&b.CreatedAt,
		&b.UpdatedAt,
		&b.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("BannerDAO.GetByID: %w", err)
	}
	return &b, nil
}

// List возвращает все не soft-deleted баннеры.
func (d *bannerDAO) List(ctx context.Context) ([]model.Banner, error) {
	rows, err := d.conn.Query(ctx, `
        SELECT id, title, content, description, created_at, updated_at, deleted_at
        FROM banners
        WHERE deleted_at IS NULL
        ORDER BY id
    `)
	if err != nil {
		return nil, fmt.Errorf("BannerDAO.List: %w", err)
	}
	defer rows.Close()

	var out []model.Banner
	for rows.Next() {
		var b model.Banner
		if err := rows.Scan(
			&b.ID,
			&b.Title,
			&b.Content,
			&b.Description,
			&b.CreatedAt,
			&b.UpdatedAt,
			&b.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("BannerDAO.List scan: %w", err)
		}
		out = append(out, b)
	}
	return out, nil
}

// Delete физически удаляет запись.
func (d *bannerDAO) Delete(ctx context.Context, id int64) error {
	cmd, err := d.conn.Exec(ctx, `
        DELETE FROM banners
        WHERE id = $1
    `, id)
	if err != nil {
		return fmt.Errorf("BannerDAO.Delete: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("BannerDAO.Delete: banner %d not found", id)
	}
	return nil
}

// SoftDelete выставляет DeletedAt = now() для метки удаления.
func (d *bannerDAO) SoftDelete(ctx context.Context, id int64) error {
	cmd, err := d.conn.Exec(ctx, `
        UPDATE banners
        SET deleted_at = $1
        WHERE id = $2 AND deleted_at IS NULL
    `, time.Now(), id)
	if err != nil {
		return fmt.Errorf("BannerDAO.SoftDelete: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("BannerDAO.SoftDelete: banner %d not found or already deleted", id)
	}
	return nil
}

// Update обновляет заголовок, контент, описание и UpdatedAt.
func (d *bannerDAO) Update(ctx context.Context, banner *model.Banner) error {
	cmd, err := d.conn.Exec(ctx, `
        UPDATE banners
        SET title       = $1,
            content     = $2,
            description = $3,
            updated_at  = $4
        WHERE id = $5 AND deleted_at IS NULL
    `, banner.Title, banner.Content, banner.Description, time.Now(), banner.ID)
	if err != nil {
		return fmt.Errorf("BannerDAO.Update: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("BannerDAO.Update: banner %d not found or deleted", banner.ID)
	}
	return nil
}
