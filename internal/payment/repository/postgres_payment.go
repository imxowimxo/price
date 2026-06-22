package repository

import (
	pay "Price/internal/payment/domain"
	"context"
	"database/sql"
	"errors"
	"time"
)

type PostPaymentRepository struct {
	db *sql.DB
}

func NewPostPaymentRepository(db *sql.DB) *PostPaymentRepository {
	return &PostPaymentRepository{db: db}
}

func (p *PostPaymentRepository) ApplyPremiumTransaction(ctx context.Context, userID int64, invoiceID string, amount int64, bankName string) error {

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now().UTC()

	future := now.AddDate(0, 1, 0)

	rows, err := tx.ExecContext(ctx, `INSERT  INTO invoice(invoice_id,bank_name,amount,user_id) VALUES ($1, $2, $3,$4) ON CONFLICT (invoice_id) DO NOTHING;`, invoiceID, bankName, amount, userID)
	if err != nil {
		return err
	}
	i, err := rows.RowsAffected()
	if err != nil {
		return err
	}
	if i == 0 {
		return pay.ErrDuplicateInvoice
	}

	rows, err = tx.ExecContext(ctx, `UPDATE users SET premium_expires_at = $1,status = $2,limit_prod = $3  WHERE id = $4`, future, "premium", 25, userID)
	if err != nil {
		return err
	}
	i, err = rows.RowsAffected()
	if err != nil {
		return err
	}
	if i == 0 {
		return errors.New("пользователь не найден в бд для оплаты подписки")
	}

	return tx.Commit()
}
