package domain

import (
	"context"
	"errors"
)

type PaymentResult struct {
	PaymentID string `json:"payment_id"`
	UserID    string `json:"user_id"`
	Status    string `json:"status"`
	Money     int    `json:"money"`
	BankName  string `json:"bank_name"`
}

type PaymentMethod interface {
	CreateInvoice(ctx context.Context, userID int64) (string, error)
	ParseCallback(ctx context.Context, res []byte, bankSign string) (*PaymentResult, error)
}

var ErrDuplicateInvoice = errors.New("invoice already exists")
