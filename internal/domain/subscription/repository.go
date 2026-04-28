package subscription

import (
	p "Price/internal/domain/product"
	"context"
)

type Repository interface {
	Create(ctx context.Context, subscription Subscription) (Subscription, error)
	GetProduct(ctx context.Context, userID int64) ([]p.Product, error)
	GetSubscribedUsers(ctx context.Context, prod int64) ([]int64, error)
	GetAll(ctx context.Context) ([]p.Product, error)
	UpdatePrice(ctx context.Context, userID int64, prodID int64, price float64) error
	Delete(ctx context.Context, userID int64, prodID int64) error
	GetSubscription(ctx context.Context, userID int64, prodID int64) (Subscription, error)
}
