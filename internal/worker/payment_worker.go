package worker

import (
	"Price/internal/domain/outbox"
	"Price/internal/domain/user"
	"context"
	"log/slog"
	"sync"
	"time"
)

type SubscriptionReminderProvider interface {
	GetUsersForReminder(ctx context.Context) ([]user.User, error)
	MarkReminderSentWithOutbox(ctx context.Context, userID int64, events outbox.SubscriptionReminderEvent) error
}
type SubscriptionNotifier struct {
	subs SubscriptionReminderProvider
	l    *slog.Logger
}

func NewSubscriptionNotifier(notifier SubscriptionReminderProvider, l *slog.Logger) *SubscriptionNotifier {
	return &SubscriptionNotifier{
		subs: notifier,
		l:    l,
	}
}

func (s *SubscriptionNotifier) Worker(ctx context.Context) error {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	s.processTick(ctx)

	for {
		select {
		case <-ctx.Done():
			s.l.Info("worker subscription reminder закончил работу")
			return nil

		case <-ticker.C:
			s.l.Info("worker subscription reminder начал обход пользователей")
			s.processTick(ctx)
		}
	}

}

func (s *SubscriptionNotifier) processTick(ctx context.Context) {

	wg := &sync.WaitGroup{}
	limit := make(chan struct{}, 10)

	list, err := s.subs.GetUsersForReminder(ctx)
	if err != nil {
		s.l.Error("ошибка получения пользователей для обхода",
			slog.Any("error", err),
		)
		return
	}

Loop:
	for _, item := range list {
		select {
		case <-ctx.Done():
			s.l.Info("обход завершен")
			break Loop
		case limit <- struct{}{}:
		}

		wg.Add(1)
		go func(user user.User) {
			defer func() { <-limit }()
			defer wg.Done()

			out := outbox.SubscriptionReminderEvent{
				UserID:    user.ID,
				ExpiresAt: user.PremiumExpiresAt,
			}

			errOut := s.subs.MarkReminderSentWithOutbox(ctx, user.ID, out)
			if errOut != nil {
				s.l.Error("ошибка отправки в outbox метод: MarkReminderSentWithOutbox",
					slog.Any("error", errOut),
				)
				return
			}
		}(item)
	}
	wg.Wait()
}

//

//

//

//

//

//

//

//

//
