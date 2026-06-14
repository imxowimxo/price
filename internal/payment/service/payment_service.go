package service

import (
	pay "Price/internal/payment/domain"
	"context"
	"errors"
	"log/slog"
	"strconv"
)

type PaymentRepository interface {
	ApplyPremiumTransaction(ctx context.Context, userID int64, invoiceID string, amount int64, bankName string) error
}

type PaymentService struct {
	repo PaymentRepository
	l    *slog.Logger
}

func NewPaymentService(repo PaymentRepository, l *slog.Logger) *PaymentService {
	return &PaymentService{repo: repo, l: l}
}

const premiumPrice = 5000

func (p *PaymentService) ProcessPayment(ctx context.Context, res *pay.PaymentResult) error {
	if res.Status != "success" {
		p.l.Info("payment service отмена платежа", "id", res.PaymentID, "userID", res.UserID)
		return nil
	}

	if res.Money != premiumPrice {
		p.l.Warn("payment service ошибка", "bad amount", res.PaymentID)
		return nil
	}

	userID, err := strconv.ParseInt(res.UserID, 10, 64)
	if err != nil {
		return err
	}

	err = p.repo.ApplyPremiumTransaction(ctx, userID, res.PaymentID, premiumPrice, res.BankName)
	if errors.Is(err, pay.ErrDuplicateInvoice) {
		return nil
	}
	if err != nil {
		p.l.Error("payment service ошибка", "id", res.PaymentID, "userID", userID, "err", err)
		return err
	}

	return nil
}
