package grpc

import (
	g "Price/gen/bot"
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) UpdateTargetPrice(ctx context.Context, req *g.UpdateRequest) (*g.UpdateResponse, error) {
	err := h.serviceSub.UpdatePrice(ctx, req.UserId, req.ProductId, req.TargetPrice)
	if err != nil {
		log.Printf("[gRPC UpdateTargetPrice] ошибка обновления цены продукта %d: %v", req.UserId, err)
		return nil, status.Error(codes.Internal, "не удалось обновить продукт, внутренняя ошибка")
	}
	return &g.UpdateResponse{
		Status: "success",
	}, nil
}
