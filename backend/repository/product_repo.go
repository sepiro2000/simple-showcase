package repository

import (
	"backend/models"
	"context"
	"database/sql"
	"fmt"
)

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	GetAll(ctx context.Context) ([]models.Product, error)
	GetByID(ctx context.Context, id int64) (models.Product, error)
	IncrementLike(ctx context.Context, id int64) error
}

type mysqlProductRepository struct {
	db *sql.DB
}

// NewMySQLProductRepository creates a new MySQL product repository
func NewMySQLProductRepository(db *sql.DB) ProductRepository {
	return &mysqlProductRepository{db: db}
}

func (r *mysqlProductRepository) GetAll(ctx context.Context) ([]models.Product, error) {
	query := `SELECT id, name, description, price, image_url, likes FROM products`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying products: %v", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.Likes); err != nil {
			return nil, fmt.Errorf("error scanning product: %v", err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %v", err)
	}

	return products, nil
}

func (r *mysqlProductRepository) GetByID(ctx context.Context, id int64) (models.Product, error) {
	query := `SELECT id, name, description, price, image_url, likes FROM products WHERE id = ?`
	var product models.Product
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.ImageURL,
		&product.Likes,
	)
	if err == sql.ErrNoRows {
		return models.Product{}, fmt.Errorf("product not found")
	}
	if err != nil {
		return models.Product{}, fmt.Errorf("error querying product: %v", err)
	}
	return product, nil
}

func (r *mysqlProductRepository) IncrementLike(ctx context.Context, id int64) error {
	query := `UPDATE products SET likes = likes + 1 WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error updating product likes: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}
