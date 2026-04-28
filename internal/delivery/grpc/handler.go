package grpc

import (
	g "Price/gen/bot"
	p "Price/internal/service/product"
	s "Price/internal/service/subscription"
	u "Price/internal/service/user"
)

type Handler struct {
	g.UnimplementedPriceServiceServer
	serviceSub     s.Service
	serviceUser    u.Service
	serviceProduct p.Service
}

func NewHandler(serviceSub s.Service, serviceUser u.Service, serviceProduct p.Service) *Handler {
	return &Handler{
		serviceSub:     serviceSub,
		serviceUser:    serviceUser,
		serviceProduct: serviceProduct,
	}
}
