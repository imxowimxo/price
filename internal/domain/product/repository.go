package product

import "context"

type Repository interface {
	Create(ctx context.Context, product Product) (Product, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, name string, newPrice float64) (Product, error)
	FindByID(ctx context.Context, id int64) (Product, error)
	UpdatePrice(ctx context.Context, prodID int64, newPrice float64) error
}
