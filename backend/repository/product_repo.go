package repository

import (
	"context"
	"database/sql"

	"backend/cache"
	"backend/models"

	"github.com/redis/go-redis/v9"
)

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	GetAll(ctx context.Context) ([]models.Product, error)
	GetByID(ctx context.Context, id int64) (*models.Product, error)
	IncrementLike(ctx context.Context, id int64) error
}

// dbProductRepository implements ProductRepository with MySQL and optional Redis
type dbProductRepository struct {
	writeDB     *sql.DB
	readDB      *sql.DB
	redisClient *redis.Client
}

// NewProductRepository creates a new instance of dbProductRepository
func NewProductRepository(writeDB, readDB *sql.DB, redisClient *redis.Client) ProductRepository {
	return &dbProductRepository{
		writeDB:     writeDB,
		readDB:      readDB,
		redisClient: redisClient,
	}
}

// GetAll retrieves all products
func (r *dbProductRepository) GetAll(ctx context.Context) ([]models.Product, error) {
	// Try to get from cache first
	if r.redisClient != nil {
		cachedProducts, exists, err := cache.GetProductList(ctx, r.redisClient)
		if err != nil {
			return nil, err
		}

		if exists {
			return cachedProducts, nil
		}
	}

	// Cache miss or Redis not available: get from database
	query := `
		SELECT id, name, description, price, image_url, likes
		FROM products
		ORDER BY id
	`

	rows, err := r.readDB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.Likes); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Cache the product list for future requests
	if r.redisClient != nil {
		if err := cache.SetProductList(ctx, r.redisClient, products); err != nil {
			// Log error but don't fail the request
			// You might want to add proper logging here
		}
	}

	return products, nil
}

// GetByID retrieves a product by its ID
func (r *dbProductRepository) GetByID(ctx context.Context, id int64) (*models.Product, error) {
	// Always get from database for real-time data
	query := `
		SELECT id, name, description, price, image_url, likes
		FROM products
		WHERE id = ?
	`

	var p models.Product
	if err := r.readDB.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.Likes,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

// IncrementLike increments the like count for a product
func (r *dbProductRepository) IncrementLike(ctx context.Context, id int64) error {
	// Update likes directly in database
	_, err := r.writeDB.ExecContext(ctx,
		"UPDATE products SET likes = likes + 1 WHERE id = ?",
		id,
	)

	// Invalidate product list cache since likes count changed
	if err == nil && r.redisClient != nil {
		if cacheErr := cache.InvalidateProductList(ctx, r.redisClient); cacheErr != nil {
			// Log cache invalidation error but don't fail the like operation
			// You might want to add proper logging here
		}
	}

	return err
}
