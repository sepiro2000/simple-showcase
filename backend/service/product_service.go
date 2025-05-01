package service

import (
	"backend/models"
	"backend/repository"
	"context"
)

// ProductService handles business logic for products
type ProductService struct {
	repo repository.ProductRepository
}

// NewProductService creates a new product service
func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// ListProducts returns all products
func (s *ProductService) ListProducts(ctx context.Context) ([]models.Product, error) {
	return s.repo.GetAll(ctx)
}

// GetProductByID returns a product by its ID
func (s *ProductService) GetProductByID(ctx context.Context, id int64) (models.Product, error) {
	return s.repo.GetByID(ctx, id)
}

// LikeProduct increments the like count for a product
func (s *ProductService) LikeProduct(ctx context.Context, id int64) error {
	return s.repo.IncrementLike(ctx, id)
}
