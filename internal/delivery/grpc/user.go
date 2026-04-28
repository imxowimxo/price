package grpc

import (
	g "Price/gen/bot"
	us "Price/internal/domain/user"
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) AddUser(ctx context.Context, req *g.CreateUser) (*g.NewUser, error) {

	user := us.User{
		ID: req.UserId,
	}
	newUser, err := h.serviceUser.Create(ctx, user)
	if err != nil {
		log.Printf("[gRPC AddUser] ошибка создания пользователя %d: %v", req.UserId, err)
		return nil, status.Error(codes.Internal, "не удалось создать пользователя, внутренняя ошибка")
	}
	result := g.NewUser{
		UserId:   newUser.ID,
		UserName: newUser.Username,
	}
	return &result, nil
}

func (h *Handler) GetUserProducts(ctx context.Context, req *g.GetUser) (*g.ProductListResponse, error) {
	products, err := h.serviceSub.ListAll(ctx, req.UserId)
	if err != nil {
		log.Printf("[gRPC GetUserProducts] ошибка получения продуктов пользователя %d: %v", req.UserId, err)
		newErr := status.Error(codes.Internal, "не удалось создать пользователя, внутренняя ошибка")
		return nil, newErr
	}
	var res g.ProductListResponse
	for _, product := range products {
		subs, err := h.serviceSub.GetSub(ctx, req.UserId, product.ID)
		if err != nil {
			log.Printf("[gRPC GetUserProducts] ошибка в получение подписки пользователя %d: %v", req.UserId, err)
			return nil, status.Error(codes.Internal, "ошибка при получении данных подписки")
		}

		res.Products = append(res.Products, &g.Product{
			ProductId:    product.ID,
			Name:         product.Name,
			Url:          product.URL,
			CurrentPrice: product.CurrentPrice,
			TargetPrice:  subs.TargetPrice,
		})
	}
	return &res, nil
}
