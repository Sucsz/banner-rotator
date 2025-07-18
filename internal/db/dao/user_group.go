package dao

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Sucsz/banner-rotator/internal/db/model"
	"github.com/jackc/pgx/v5"
)

// UserGroupDAO — интерфейс для таблицы user_groups
type UserGroupDAO interface {
	Create(ctx context.Context, group *model.UserGroup) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.UserGroup, error)
	List(ctx context.Context) ([]model.UserGroup, error)
	Delete(ctx context.Context, id int64) error
	SoftDelete(ctx context.Context, id int64) error
	Update(ctx context.Context, group *model.UserGroup) error
}

type userGroupDAO struct {
	conn *pgx.Conn
}

// NewUserGroupDAO создаёт экземпляр userGroupDAO в виде интерфейса UserGroupDAO
func NewUserGroupDAO(conn *pgx.Conn) UserGroupDAO {
	return &userGroupDAO{conn: conn}
}

// Create вставляет новую запись и возвращает её ID
func (d *userGroupDAO) Create(ctx context.Context, group *model.UserGroup) (int64, error) {
	var id int64
	now := time.Now()
	err := d.conn.QueryRow(ctx, `
        INSERT INTO user_groups (description, created_at, updated_at)
        VALUES ($1, $2, $3)
        RETURNING id
    `, group.Description, now, now).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("UserGroupDAO.Create: %w", err)
	}
	return id, nil
}

// GetByID возвращает пользовательскую группу по ID, исключая soft-deleted
func (d *userGroupDAO) GetByID(ctx context.Context, id int64) (*model.UserGroup, error) {
	row := d.conn.QueryRow(ctx, `
        SELECT id, description, created_at, updated_at, deleted_at
        FROM user_groups
        WHERE id = $1 AND deleted_at IS NULL
    `, id)

	var g model.UserGroup
	err := row.Scan(
		&g.ID,
		&g.Description,
		&g.CreatedAt,
		&g.UpdatedAt,
		&g.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("UserGroupDAO.GetByID: %w", err)
	}
	return &g, nil
}

// List возвращает все не soft-deleted пользовательские группы
func (d *userGroupDAO) List(ctx context.Context) ([]model.UserGroup, error) {
	rows, err := d.conn.Query(ctx, `
        SELECT id, description, created_at, updated_at, deleted_at
        FROM user_groups
        WHERE deleted_at IS NULL
        ORDER BY id
    `)
	if err != nil {
		return nil, fmt.Errorf("UserGroupDAO.List: %w", err)
	}
	defer rows.Close()

	var out []model.UserGroup
	for rows.Next() {
		var g model.UserGroup
		if err := rows.Scan(
			&g.ID,
			&g.Description,
			&g.CreatedAt,
			&g.UpdatedAt,
			&g.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("UserGroupDAO.List scan: %w", err)
		}
		out = append(out, g)
	}
	return out, nil
}

// Delete физически удаляет запись
func (d *userGroupDAO) Delete(ctx context.Context, id int64) error {
	cmd, err := d.conn.Exec(ctx, `DELETE FROM user_groups WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("UserGroupDAO.Delete: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("UserGroupDAO.Delete: user_group %d not found", id)
	}
	return nil
}

// SoftDelete выставляет DeletedAt = now() для метки удаления
func (d *userGroupDAO) SoftDelete(ctx context.Context, id int64) error {
	cmd, err := d.conn.Exec(ctx, `
        UPDATE user_groups
        SET deleted_at = $1
        WHERE id = $2 AND deleted_at IS NULL
    `, time.Now(), id)
	if err != nil {
		return fmt.Errorf("UserGroupDAO.SoftDelete: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("UserGroupDAO.SoftDelete: user_group %d not found or already deleted", id)
	}
	return nil
}

// Update обновляет описание и UpdatedAt
func (d *userGroupDAO) Update(ctx context.Context, group *model.UserGroup) error {
	cmd, err := d.conn.Exec(ctx, `
        UPDATE user_groups
        SET description = $1, updated_at = $2
        WHERE id = $3 AND deleted_at IS NULL
    `, group.Description, time.Now(), group.ID)
	if err != nil {
		return fmt.Errorf("UserGroupDAO.Update: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("UserGroupDAO.Update: user_group %d not found or deleted", group.ID)
	}
	return nil
}
