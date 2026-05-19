package worker

import (
	par "Price/gen/parser"
	"Price/internal/domain/price_drop_event"
	"Price/internal/domain/product"
	inf "Price/internal/infrastructure"
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"
)

type PriceFetcher interface {
	FetchPrice(ctx context.Context, url string) (float64, error)
}

type ProductProvider interface {
	GetProductsToCheck(ctx context.Context) ([]product.Product, error)
	UpdatePriceWithOutbox(ctx context.Context, productID int64, newPrice float64, events []price_drop_event.PriceDropEvent) error
}

//type Outbox interface {
//	UpdatePriceWithOutbox(ctx context.Context, productID int64, newPrice float64, events []price_drop_event.PriceDropEvent) error
//}

type SubscriptionProvider interface {
	GetUsersForPriceDrop(ctx context.Context, productID int64, currentPrice float64) ([]int64, error)
}

type PriceWatcher struct {
	productProvider ProductProvider
	//outbox          Outbox
	subProvider SubscriptionProvider
	l           *slog.Logger
	parser      par.GetPriceClient
	br          *inf.CircuitBreaker
}

func NewPriceWatcher(pp ProductProvider, sp SubscriptionProvider, l *slog.Logger, parser par.GetPriceClient, br *inf.CircuitBreaker) *PriceWatcher {
	return &PriceWatcher{
		productProvider: pp,
		//outbox:          outbox,
		subProvider: sp,
		l:           l,
		parser:      parser,
		br:          br,
	}
}

func (w *PriceWatcher) Worker(ctx context.Context) error {

	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	w.processTick(ctx)

	for {
		select {
		case <-ctx.Done():
			w.l.Info("PriceWatcher Worker завершился ")
			return nil

		case <-ticker.C:
			w.l.Info("горутина начала обход")
			w.processTick(ctx)
		}
	}

}

func (w *PriceWatcher) processTick(ctx context.Context) {

	wg := &sync.WaitGroup{}
	limit := make(chan struct{}, 10)

	list, err := w.productProvider.GetProductsToCheck(ctx)
	if err != nil {
		w.l.Error("ошибка получения продуктов для обхода",
			slog.Any("error", err),
		)
		return
	}
Loop:
	for _, item := range list {

		select {
		case <-ctx.Done():
			w.l.Info("обход завершен")
			break Loop
		case limit <- struct{}{}:
		}

		wg.Add(1)
		go func(p product.Product) {
			defer func() { <-limit }()
			defer wg.Done()

			parserReq := par.ParseRequest{
				Url: p.URL,
			}

			var price *par.ParseResponse
			var err error

			for attempt := 1; attempt <= 3; attempt++ {
				ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)

				err = w.br.Execute(func() error {
					var grpcErr error
					price, grpcErr = w.parser.ParserWeb(ctxTimeout, &parserReq)
					return grpcErr
				})
				cancel()

				if err == nil {
					break
				}
				w.l.Warn("не удалось получить цену, пробуем еще...", slog.Int64("попытка", int64(attempt)), slog.Any("error", err))
				if attempt < 3 {
					select {
					case <-ctx.Done():
						return
					case <-time.After(2 * time.Second):
					}
				}
			}

			if err != nil {
				if errors.Is(err, inf.ErrCircuitOpen) {
					w.l.Warn("не удалось получить цену, пробуем еще...",
						slog.Any("error", err),
					)
					return
				}

				w.l.Error("ошибка получения цены",
					slog.String("url", p.URL),
					slog.Int64("product_id", p.ID),
					slog.Any("error", err),
				)
				return
			}

			if price.Price != p.CurrentPrice {
				var events []price_drop_event.PriceDropEvent

				if price.Price < p.CurrentPrice && p.CurrentPrice > 0 {
					idUser, err := w.subProvider.GetUsersForPriceDrop(ctx, p.ID, price.Price)
					if err != nil {
						w.l.Error("ошибка получения пользователя для отправки уведомления",
							slog.String("url", p.URL),
							slog.Int64("product_id", p.ID),
							slog.Any("error", err),
						)
						return
					}

					for _, user := range idUser {
						event := price_drop_event.PriceDropEvent{
							UserID:    user,
							ProductID: p.ID,
							OldPrice:  p.CurrentPrice,
							NewPrice:  price.Price,
						}
						events = append(events, event)
					}
				}

				err = w.productProvider.UpdatePriceWithOutbox(ctx, p.ID, price.Price, events)
				if err != nil {
					w.l.Error("ошибка функции UpdatePriceWithOutbox", slog.Any("error", err))
					return
				}
			}

		}(item)
	}
	wg.Wait()
	w.l.Info("обход завершен", slog.Int("processed_count", len(list)))
}
