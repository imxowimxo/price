package worker

import (
	r "Price/internal/repository"
	"context"
	"log"
	"time"
)

type Notifier interface {
	SendMessage(ctx context.Context, topic string, key string, payload []byte) error
}

type OutboxRelay struct {
	repo   *r.OutboxRepo
	n      Notifier
	logger *log.Logger
}

func NewOutboxRelay(repo *r.OutboxRepo, n Notifier, logger *log.Logger) *OutboxRelay {
	return &OutboxRelay{
		repo:   repo,
		n:      n,
		logger: logger,
	}
}

func (r *OutboxRelay) Run(ctx context.Context) error {

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			events, err := r.repo.GetPendingEvents(ctx, 100)
			if err != nil {
				r.logger.Println(err)
			}
			if len(events) == 0 {
				continue
			}
			for _, event := range events {

				err = r.n.SendMessage(ctx, event.Topic, event.MessageKey, event.Payload)
				if err != nil {
					r.logger.Println(err)
					continue
				}

				err = r.repo.MarkEventAsSent(ctx, event.ID)
				if err != nil {
					r.logger.Println(err)
				}
			}
		}
	}
}
