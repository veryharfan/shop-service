package db

import (
	"context"
	"database/sql"
	"log/slog"
	"shop-service/app/domain"
)

type shopRepository struct {
	conn *sql.DB
}

func NewShopRepository(db *sql.DB) domain.ShopRepository {
	return &shopRepository{db}
}

func (r *shopRepository) Create(ctx context.Context, req *domain.Shop, tx *sql.Tx) error {
	query := `INSERT INTO shops (user_id, name) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	row := tx.QueryRowContext(ctx, query, req.UserID, req.Name)
	if err := row.Scan(&req.ID, &req.CreatedAt, &req.UpdatedAt); err != nil {
		slog.ErrorContext(ctx, "[shopRepository] Create", "scan", err)
		return err
	}

	return nil
}

func (r *shopRepository) GetByUserID(ctx context.Context, userID int64) (*domain.Shop, error) {
	query := `SELECT id, user_id, name, created_at, updated_at FROM shops WHERE user_id = $1`
	row := r.conn.QueryRowContext(ctx, query, userID)

	shop := &domain.Shop{}
	if err := row.Scan(&shop.ID, &shop.UserID, &shop.Name, &shop.CreatedAt, &shop.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		slog.ErrorContext(ctx, "[shopRepository] GetByUserID", "scan", err)
		return nil, err
	}

	return shop, nil
}

func (r *shopRepository) BeginTransaction(ctx context.Context) (*sql.Tx, error) {
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		slog.ErrorContext(ctx, "[shopRepository] BeginTransaction", "beginTx", err)
		return nil, err
	}
	return tx, nil
}

func (r *shopRepository) WithTransaction(ctx context.Context, tx *sql.Tx, fn func(context.Context, *sql.Tx) error) error {
	if err := fn(ctx, tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			slog.ErrorContext(ctx, "[shopRepository] WithTransaction", "rollback", rollbackErr)
			return rollbackErr
		}
		return err
	}
	return tx.Commit()
}
