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
	query := `
		SELECT id, name, description, price, image_url
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
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL); err != nil {
			return nil, err
		}

		// Get likes from Redis if available
		if r.redisClient != nil {
			likes, err := cache.GetProductLikes(ctx, r.redisClient, p.ID)
			if err != nil {
				return nil, err
			}
			p.Likes = likes
		} else {
			// Get likes from database if Redis is not available
			if err := r.readDB.QueryRowContext(ctx, "SELECT likes FROM products WHERE id = ?", p.ID).Scan(&p.Likes); err != nil {
				return nil, err
			}
		}

		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// GetByID retrieves a product by its ID
func (r *dbProductRepository) GetByID(ctx context.Context, id int64) (*models.Product, error) {
	query := `
		SELECT id, name, description, price, image_url
		FROM products
		WHERE id = ?
	`

	var p models.Product
	if err := r.readDB.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Get likes from Redis if available
	if r.redisClient != nil {
		likes, err := cache.GetProductLikes(ctx, r.redisClient, p.ID)
		if err != nil {
			return nil, err
		}
		p.Likes = likes
	} else {
		// Get likes from database if Redis is not available
		if err := r.readDB.QueryRowContext(ctx, "SELECT likes FROM products WHERE id = ?", p.ID).Scan(&p.Likes); err != nil {
			return nil, err
		}
	}

	return &p, nil
}

// IncrementLike increments the like count for a product
func (r *dbProductRepository) IncrementLike(ctx context.Context, id int64) error {
	if r.redisClient != nil {
		// Increment likes in Redis
		newLikes, err := cache.IncrementProductLikes(ctx, r.redisClient, id)
		if err != nil {
			return err
		}

		// Update the database with the new like count
		_, err = r.writeDB.ExecContext(ctx,
			"UPDATE products SET likes = ? WHERE id = ?",
			newLikes, id,
		)
		return err
	}

	// If Redis is not available, increment likes directly in the database
	_, err := r.writeDB.ExecContext(ctx,
		"UPDATE products SET likes = likes + 1 WHERE id = ?",
		id,
	)
	return err
}
