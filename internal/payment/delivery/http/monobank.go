package http

import (
	p "Price/internal/payment/domain"
	"context"
	"io"
	"net/http"
)

type MonoBankInterface interface {
	CreateInvoice(ctx context.Context, userID int64) (string, error)
	ParseCallback(ctx context.Context, res []byte, bankSign string) (*p.PaymentResult, error)
}

type PaymentServiceInterface interface {
	ProcessPayment(ctx context.Context, res *p.PaymentResult) error
}

type Handler struct {
	PaymentServiceInterface PaymentServiceInterface
	invoiceParser           MonoBankInterface
}

func NewHandler(PaymentServiceInterface PaymentServiceInterface, invoiceParser MonoBankInterface) *Handler {
	return &Handler{
		PaymentServiceInterface: PaymentServiceInterface,
		invoiceParser:           invoiceParser,
	}
}

func (h *Handler) ParserWebhook(w http.ResponseWriter, r *http.Request) {

	xSign := r.Header.Get("X-Sign")
	if xSign == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Не удалось прочитать тело запроса", http.StatusBadRequest)
		return
	}

	res, err := h.invoiceParser.ParseCallback(r.Context(), bodyBytes, xSign)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.PaymentServiceInterface.ProcessPayment(r.Context(), res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return

}
