package repository

import (
	p "Price/internal/domain/product"
	s "Price/internal/domain/subscription"
	"context"
	"database/sql"
	"errors"
)

type PostgresSubscriptionRepository struct {
	db *sql.DB
}

func NewPostgresSubscriptionRepository(db *sql.DB) *PostgresSubscriptionRepository {
	return &PostgresSubscriptionRepository{db: db}
}

func (ps *PostgresSubscriptionRepository) Create(ctx context.Context, subscription s.Subscription) (s.Subscription, error) {

	query := `INSERT INTO subscriptions ( user_id, product_id, target_price,is_triggered) VALUES ($1, $2, $3,$4)`

	_, err := ps.db.ExecContext(ctx, query, subscription.UserID, subscription.ProductID, subscription.TargetPrice, subscription.IsTriggered)
	if err != nil {
		return s.Subscription{}, err
	}
	return subscription, nil
}

func (ps *PostgresSubscriptionRepository) GetProduct(ctx context.Context, userID int64) ([]p.Product, error) {

	query := `SELECT products.id,products.url,products.price,products.name FROM subscriptions JOIN products ON subscriptions.product_id = products.id WHERE subscriptions.user_id = $1`
	rows, err := ps.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]p.Product, 0)
	for rows.Next() {
		product := p.Product{}
		err1 := rows.Scan(&product.ID, &product.URL, &product.CurrentPrice, &product.Name)
		if err1 != nil {
			return nil, err1
		}
		list = append(list, product)
	}

	if err2 := rows.Err(); err2 != nil {
		return nil, err2
	}

	return list, nil
}

func (ps *PostgresSubscriptionRepository) GetSubscribedUsers(ctx context.Context, prod int64) ([]int64, error) {

	query := `SELECT user_id FROM subscriptions WHERE product_id = $1`
	rows, err := ps.db.QueryContext(ctx, query, prod)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]int64, 0)
	for rows.Next() {
		var userID int64
		err1 := rows.Scan(&userID)

		if err1 != nil {
			return nil, err1
		}
		list = append(list, userID)
	}

	if err2 := rows.Err(); err2 != nil {
		return nil, err2
	}
	return list, nil
}

func (ps *PostgresSubscriptionRepository) UpdatePrice(ctx context.Context, userID int64, prodID int64, price float64) error {
	_, err := ps.db.ExecContext(ctx, "UPDATE products SET price = $1 WHERE user_id = $2 AND product_id = $3", price, userID, prodID)
	return err
}

func (ps *PostgresSubscriptionRepository) GetAll(ctx context.Context) ([]p.Product, error) {
	query := `SELECT id, url, price FROM products`
	rows, err := ps.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	list := make([]p.Product, 0)
	for rows.Next() {
		var prod p.Product
		err1 := rows.Scan(&prod.ID, &prod.URL, &prod.CurrentPrice)
		if err1 != nil {
			return nil, err1
		}
		list = append(list, prod)

	}
	if err2 := rows.Err(); err2 != nil {
		return nil, err2
	}
	return list, nil
}

func (ps *PostgresSubscriptionRepository) Delete(ctx context.Context, userID int64, prodID int64) error {
	res, err := ps.db.ExecContext(ctx, "DELETE FROM subscriptions WHERE user_id = $1 AND product_id = $2", userID, prodID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("не найдена подписка для удаления")
	}
	return nil
}

func (ps *PostgresSubscriptionRepository) GetSubscription(ctx context.Context, userID int64, prodID int64) (s.Subscription, error) {
	query := `SELECT id, user_id, product_id, target_price, is_triggered FROM subscriptions WHERE user_id = $1 AND product_id = $2`
	var sub s.Subscription
	err := ps.db.QueryRowContext(ctx, query, userID, prodID).Scan(&sub.ID, &sub.UserID, &sub.ProductID, &sub.TargetPrice, &sub.IsTriggered)

	if errors.Is(err, sql.ErrNoRows) {
		return s.Subscription{}, errors.New("у пользователя нет продуктов в подписке")
	}

	if err != nil {
		return s.Subscription{}, err
	}

	return sub, nil
}

func (ps *PostgresSubscriptionRepository) GetUsersForPriceDrop(ctx context.Context, productID int64, currentPrice float64) ([]int64, error) {

	query := `SELECT user_id FROM subscriptions WHERE product_id = $1 AND target_price >= $2 AND is_triggered = false`

	rows, err := ps.db.QueryContext(ctx, query, productID, currentPrice)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]int64, 0)

	for rows.Next() {
		var userID int64

		err1 := rows.Scan(&userID)
		if err1 != nil {
			return nil, err1
		}

		list = append(list, userID)
	}

	if err2 := rows.Err(); err2 != nil {
		return nil, err2
	}

	return list, nil
}
