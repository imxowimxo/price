package repository

import (
	us "Price/internal/domain/user"
	"context"
	"database/sql"
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
    // бот спамил ошибками при нажатии кнопок, потому что пытался
    // создать юзера, который уже есть, и ловил ошибку уникальности tg_id.
	query := `INSERT INTO users(username,tg_id) VALUES ($1, $2) RETURNING id`
	err := p.db.QueryRowContext(ctx, query, user.Username, user.TgID).Scan(&user.ID)
	if err != nil {
		return us.User{}, err
	}
	return user, nil
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
