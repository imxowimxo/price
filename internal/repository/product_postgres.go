package repository

import (
	"Price/internal/domain/price_drop_event"
	pr "Price/internal/domain/product"
	"context"
	"database/sql"
	"encoding/json"
)

type PostgresProductRepo struct {
	db *sql.DB
}

func NewPostgresProductRepo(db *sql.DB) *PostgresProductRepo {
	return &PostgresProductRepo{db: db}
}

func (ps *PostgresProductRepo) GetProductsToCheck(ctx context.Context) ([]pr.Product, error) {

	query := `SELECT id,url,current_price,name FROM products`
	rows, err := ps.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var list []pr.Product
	for rows.Next() {
		product := pr.Product{}
		err = rows.Scan(&product.ID, &product.URL, &product.CurrentPrice, &product.Name)
		if err != nil {
			return nil, err
		}

		list = append(list, product)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (ps *PostgresProductRepo) Create(ctx context.Context, product pr.Product) (pr.Product, error) {
	query := `INSERT INTO products (url, current_price,name) VALUES ($1, $2, $3) RETURNING id`
	err := ps.db.QueryRowContext(ctx, query, product.URL, product.CurrentPrice, product.Name).Scan(&product.ID)

	if err != nil {
		return pr.Product{}, err
	}

	return product, nil
}

func (ps *PostgresProductRepo) Delete(ctx context.Context, id int64) error {

	_, err := ps.db.ExecContext(ctx, "DELETE FROM products WHERE id = $1", id)

	if err != nil {
		return err
	}
	return nil
}

func (ps *PostgresProductRepo) Update(ctx context.Context, id int64, name string, newPrice float64) (pr.Product, error) {
	query := `
        UPDATE products 
        SET current_price = $1, name = $2 
        WHERE id = $3 
        RETURNING id, url, current_price, name
    `

	var product pr.Product

	err := ps.db.QueryRowContext(ctx, query, newPrice, name, id).Scan(
		&product.ID,
		&product.URL,
		&product.CurrentPrice,
		&product.Name,
	)

	if err != nil {
		return pr.Product{}, err
	}

	return product, nil
}

func (ps *PostgresProductRepo) FindByID(ctx context.Context, id int64) (pr.Product, error) {
	query := `SELECT id, url, current_price, name FROM products WHERE id = $1`

	product := pr.Product{}

	err := ps.db.QueryRowContext(ctx, query, id).Scan(&product.ID, &product.URL, &product.CurrentPrice, &product.Name)
	if err != nil {
		return pr.Product{}, err
	}
	return product, nil
}

func (ps *PostgresProductRepo) UpdatePrice(ctx context.Context, prodID int64, newPrice float64) error {
	query := `UPDATE products SET current_price = $1 WHERE id = $2`
	_, err := ps.db.ExecContext(ctx, query, newPrice, prodID)
	return err
}

func (ps *PostgresProductRepo) UpdatePriceWithOutbox(ctx context.Context, prodID int64, newPrice float64, events []price_drop_event.PriceDropEvent) error {

	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `UPDATE products SET current_price = $1 WHERE id = $2`, newPrice, prodID)
	if err != nil {
		return err
	}

	for _, event := range events {
		jsonBytes, err := json.Marshal(&event)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, "INSERT INTO outbox_events (payload) VALUES ($1)", jsonBytes)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	return err
}

//

//

//

//

//

//

//
