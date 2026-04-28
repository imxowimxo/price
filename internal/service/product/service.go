package product

import (
	"Price/internal/domain/product"
	"context"
	"errors"
)

type Service interface {
	Create(ctx context.Context, product product.Product) (product.Product, error)
	Update(ctx context.Context, p product.Product) (product.Product, error)
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (product.Product, error)
	UpdatePrice(ctx context.Context, productID int64, newPrice float64) error
}

type service struct {
	repo product.Repository
}

func NewService(repo product.Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, newProduct product.Product) (product.Product, error) {

	if newProduct.Name == "" {
		return newProduct, errors.New("имя продукта не может быть пустым")
	}
	if newProduct.URL == "" {
		return newProduct, errors.New("ссылка на продукт не может быть пустой")
	}

	createdProduct, err := s.repo.Create(ctx, newProduct)
	if err != nil {
		return product.Product{}, err
	}
	return createdProduct, nil

}

func (s *service) Update(ctx context.Context, p product.Product) (product.Product, error) {
	if p.ID == 0 {
		return product.Product{}, errors.New("пустое айди")
	}
	if p.Name == "" {
		return product.Product{}, errors.New("имя не может быть пустым")
	}
	if p.CurrentPrice <= 0 {
		return product.Product{}, errors.New("цена продукта должна быть больше нуля")
	}

	return s.repo.Update(ctx, p.ID, p.Name, p.CurrentPrice)
}

func (s *service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) FindByID(ctx context.Context, id int64) (product.Product, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) UpdatePrice(ctx context.Context, productID int64, newPrice float64) error {
	if productID <= 0 {
		return errors.New("пустое айди")
	}
	if newPrice <= 0 {
		return errors.New("цена продукта должна быть больше нуля")
	}

	return s.repo.UpdatePrice(ctx, productID, newPrice)
}
