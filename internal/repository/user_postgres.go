package repository

import (
	us "Price/internal/domain/user"
	"context"
	"database/sql"
	"errors"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (p *PostgresRepository) Create(ctx context.Context, user us.User) (us.User, error) {

	query := `INSERT INTO users(username,tg_id) VALUES ($1, $2) ON CONFLICT (tg_id) DO NOTHING RETURNING id`

	err := p.db.QueryRowContext(ctx, query, user.Username, user.TgID).Scan(&user.ID)

	if err == nil {
		return user, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		query = `SELECT id FROM users WHERE tg_id = $1`
		err = p.db.QueryRowContext(ctx, query, user.TgID).Scan(&user.ID)
		if err == nil {
			return user, nil
		}
	}
	return us.User{}, err
}

func (p *PostgresRepository) GetByID(ctx context.Context, userID int64) (us.User, error) {
	query := `SELECT id,username,tg_id FROM users WHERE id = $1`
	row := p.db.QueryRowContext(ctx, query, userID)
	var user us.User
	err := row.Scan(&user.ID, &user.Username, &user.TgID)
	if err != nil {
		return us.User{}, err
	}
	return user, nil
}
