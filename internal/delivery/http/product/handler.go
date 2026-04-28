package product

import (
	"Price/internal/domain/product"
	sv "Price/internal/service/product"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type HTTPHandler struct {
	service sv.Service
}
type updateProductRequest struct {
	Name     string  `json:"name"`
	NewPrice float64 `json:"new_price"`
}

type updateProductResponse struct {
	Success bool `json:"success"`
}

func NewHTTPHandler(service sv.Service) *HTTPHandler {
	return &HTTPHandler{service: service}
}

func (h *HTTPHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	req := updateProductRequest{}

	idStr := chi.URLParam(r, "id")

	ctx := r.Context()
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	prod := product.Product{
		ID:           newID,
		Name:         req.Name,
		CurrentPrice: req.NewPrice,
	}

	_, err = h.service.Update(ctx, prod)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(updateProductResponse{Success: true})
}

type createProductRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"new_price"`
	URL   string  `json:"url"`
}

type createProductResponse struct {
	ID int64 `json:"id"`
}

func (h *HTTPHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	req := createProductRequest{}
	ctx := r.Context()
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	prod := product.Product{
		Name:         req.Name,
		CurrentPrice: req.Price,
		URL:          req.URL,
	}

	createdProd, err := h.service.Create(ctx, prod)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createProductResponse{ID: createdProd.ID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *HTTPHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")
	ctx := r.Context()

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "неверный формат ID", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}

type getProductsResponse struct {
	ID           int64   `json:"id"`
	URL          string  `json:"url"`
	CurrentPrice float64 `json:"current_price"`
	Name         string  `json:"name"`
}

func (h *HTTPHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	ctx := r.Context()

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "неверный формат ID", http.StatusBadRequest)
		return
	}

	res, err := h.service.FindByID(ctx, id)
	if err != nil {
		http.Error(w, "продукт не найден", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(getProductsResponse{
		ID:           res.ID,
		URL:          res.URL,
		CurrentPrice: res.CurrentPrice,
		Name:         res.Name,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
