package subscription

import (
	"Price/internal/domain/product"
	sub "Price/internal/domain/subscription"
	svs "Price/internal/service/subscription"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type HTTPHandler struct {
	service svs.Service
}

type subscriptionCreateRequest struct {
	UserID      int64   `json:"user_id"`
	ProductID   int64   `json:"product_id"`
	TargetPrice float64 `json:"target_price"`
}

type subscriptionCreateResponse struct {
	ID int64 `json:"id"`
}

func NewHTTPHandler(s svs.Service) *HTTPHandler {
	return &HTTPHandler{service: s}
}

func (h *HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := subscriptionCreateRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	subs := sub.Subscription{
		UserID:      req.UserID,
		ProductID:   req.ProductID,
		TargetPrice: req.TargetPrice,
	}

	res, err := h.service.Create(ctx, subs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(subscriptionCreateResponse{ID: res.ID})
}

type subscriptionListResponse struct {
	List []product.Product `json:"list"`
}

func (h *HTTPHandler) List(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")
	ctx := r.Context()

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "неверный формат ID", http.StatusBadRequest)
		return
	}

	result, err := h.service.ListAll(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(subscriptionListResponse{List: result})

}
