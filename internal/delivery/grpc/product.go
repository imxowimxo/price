package grpc

import (
	g "Price/gen/bot"
	"Price/internal/domain/product"
	"Price/internal/domain/subscription"
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) AddProduct(ctx context.Context, req *g.CreateProductRequest) (*g.CreateProductResponse, error) {

	prod := product.Product{
		URL:          req.Url,
		CurrentPrice: 0,
		Name:         req.Name,
	}

	res, err := h.serviceProduct.Create(ctx, prod)
	if err != nil {
		log.Printf("[gRPC AddProduct] ошибка создания продукта %s: %v", req.Name, err)
		return nil, status.Error(codes.Internal, "не удалось создать продукт, внутренняя ошибка")
	}

	sub := subscription.Subscription{
		UserID:      req.UserId,
		ProductID:   res.ID,
		TargetPrice: req.TargetPrice,
		IsTriggered: false,
	}
	_, err = h.serviceSub.Create(ctx, sub)
	if err != nil {
		log.Printf("[gRPC AddProduct] ошибка создания подписки: %v", err)
		return nil, status.Error(codes.Internal, "не удалось создать подписку, внутренняя ошибка")
	}

	newproduct := &g.CreateProductResponse{
		ProductId: res.ID,
	}
	return newproduct, nil
}

func (h *Handler) Delete(ctx context.Context, req *g.DeleteProductRequest) (*g.DeleteProductResponse, error) {
	err := h.serviceSub.Delete(ctx, req.UserId, req.ProductId)
	if err != nil {
		log.Printf("[gRPC Delete] ошибка при удалении подписки: %v", err)
		return nil, status.Error(codes.Internal, "не удалось удалить подписку, внутренняя ошибка")
	}
	stat := g.DeleteProductResponse{Status: "success"}
	return &stat, nil
}

func (h *Handler) Get(ctx context.Context, req *g.GetProductRequest) (*g.Product, error) {
	prod, err := h.serviceProduct.FindByID(ctx, req.ProductId)
	if err != nil {
		log.Printf("[gRPC Get] ошибка при попытке найти продукт: %v", err)
		return nil, status.Error(codes.Internal, "не удалось найти продукт, внутренняя ошибка")
	}
	sub, err := h.serviceSub.GetSub(ctx, req.UserId, req.ProductId)
	if err != nil {
		log.Printf("[gRPC Get] ошибка поиска подписки: %v", err)
		return nil, status.Error(codes.Internal, "не удалось найти данные подписки")
	}

	res := &g.Product{
		Name:         prod.Name,
		Url:          prod.URL,
		CurrentPrice: prod.CurrentPrice,
		TargetPrice:  sub.TargetPrice,
		ProductId:    prod.ID,
	}
	return res, nil
}
