package worker

import (
	"Price/internal/domain/price_drop_event"
	r "Price/internal/repository"
	"context"
	"encoding/json"
	"log"
	"time"
)

type Notifier interface {
	SendPriceDrop(ctx context.Context, event price_drop_event.PriceDropEvent) error
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
				var priceEvent price_drop_event.PriceDropEvent
				if err := json.Unmarshal(event.Payload, &priceEvent); err != nil {
					r.logger.Println(err)
				}
				err = r.n.SendPriceDrop(ctx, priceEvent)
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
