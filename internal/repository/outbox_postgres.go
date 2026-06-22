package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type OutboxRepo struct {
	db *sql.DB
}
type OutboxEvent struct {
	ID         uuid.UUID
	Payload    []byte
	Status     string
	CreatedAt  time.Time
	Topic      string
	MessageKey string
}

func NewOutboxRepo(db *sql.DB) *OutboxRepo {
	return &OutboxRepo{db: db}
}

func (o *OutboxRepo) GetPendingEvents(ctx context.Context, limit int) ([]OutboxEvent, error) {

	var events []OutboxEvent

	query := `SELECT id, payload,topic_name,message_key FROM outbox_events WHERE status = 'pending' ORDER BY created_at ASC LIMIT $1`

	rows, err := o.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event OutboxEvent
		if err := rows.Scan(&event.ID, &event.Payload, &event.Topic, &event.MessageKey); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func (o *OutboxRepo) MarkEventAsSent(ctx context.Context, eventID uuid.UUID) error {
	_, err := o.db.ExecContext(ctx, `UPDATE outbox_events SET status = 'sent' WHERE id = $1`, eventID)
	if err != nil {
		return err
	}
	return nil
}
