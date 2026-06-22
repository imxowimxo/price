package repository

import (
	"Price/internal/domain/outbox"
	out "Price/internal/domain/outbox"
	us "Price/internal/domain/user"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
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

func (p *PostgresRepository) MarkReminderSentWithOutbox(ctx context.Context, userID int64, event outbox.SubscriptionReminderEvent) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	payloadBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	strUserID := fmt.Sprintf("%d", userID)

	query := `
        UPDATE users 
        SET reminder = true
        WHERE id = $1 
    `

	_, err = tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO outbox_events (payload, topic_name, message_key) VALUES ($1, $2, $3)`,
		payloadBytes, out.TopicSubscriptions, strUserID,
	)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func (p *PostgresRepository) GetUsersForReminder(ctx context.Context) ([]us.User, error) {

	day := time.Now().UTC()
	day = day.AddDate(0, 0, 2)

	query := `SELECT id,username,tg_id,premium_expires_at FROM users WHERE status = 'premium' AND reminder = false AND premium_expires_at <= $1`
	rows, err := p.db.QueryContext(ctx, query, day)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []us.User
	for rows.Next() {
		var user us.User
		err = rows.Scan(&user.ID, &user.Username, &user.TgID, &user.PremiumExpiresAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
