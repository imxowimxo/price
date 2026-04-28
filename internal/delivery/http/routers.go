package http

import (
	HP "Price/internal/delivery/http/product"
	SB "Price/internal/delivery/http/subscription"

	"github.com/go-chi/chi/v5"
)

func NewRouter(prod *HP.HTTPHandler, subs *SB.HTTPHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/products", prod.CreateProduct)
	r.Get("/products/{id}", prod.GetProduct)
	r.Put("/products/{id}", prod.UpdateProduct)
	r.Delete("/products/{id}", prod.DeleteProduct)

	r.Post("/subscription", subs.Create)
	r.Get("/subscription/{id}", subs.List)

	return r
}
