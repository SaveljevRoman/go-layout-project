package service

import (
	"context"
	"github.com/SaveljevRoman/go-layout-project/internal/models"
	"log"
	"time"
)

type ProductRepository interface {
	GetByID(ctx context.Context, id int64) (*models.Product, error)
	GetAll(ctx context.Context) ([]*models.Product, error)
	Create(ctx context.Context, product *models.Product) (int64, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id int64) error
}

type ProductCache interface {
	GetByID(ctx context.Context, id int64) (*models.Product, error)
	Set(ctx context.Context, product *models.Product, expiration time.Duration) error
	Delete(ctx context.Context, id int64) error
	SetAllProducts(ctx context.Context, products []*models.Product, expiration time.Duration) error
}

type ProductService struct {
	repo  ProductRepository
	cache ProductCache
}

func NewProductService(repo ProductRepository, cache ProductCache) *ProductService {
	return &ProductService{
		repo:  repo,
		cache: cache,
	}
}

func (s *ProductService) GetProduct(ctx context.Context, id int64) (*models.Product, error) {
	// Сначала пытаемся получить из кеша
	product, err := s.cache.GetByID(ctx, id)
	if err != nil {
		log.Printf("Cache error: %v", err)
	}

	if product != nil {
		return product, nil
	}

	// Если в кеше нет, получаем из БД
	product, err = s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if product != nil {
		// Кешируем результат на 5 минут
		if err := s.cache.Set(ctx, product, 5*time.Minute); err != nil {
			log.Printf("Failed to cache product: %v", err)
		}
	}

	return product, nil
}

func (s *ProductService) GetAllProducts(ctx context.Context) ([]*models.Product, error) {
	return s.repo.GetAll(ctx)
}

func (s *ProductService) CreateProduct(ctx context.Context, product *models.Product) (int64, error) {
	id, err := s.repo.Create(ctx, product)
	if err != nil {
		return 0, err
	}

	// Обновить продукт с ID
	product.ID = id
	if err := s.cache.Set(ctx, product, 5*time.Minute); err != nil {
		log.Printf("Failed to cache product: %v", err)
	}

	return id, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, product *models.Product) error {
	if err := s.repo.Update(ctx, product); err != nil {
		return err
	}

	// Обновляем кеш
	if err := s.cache.Set(ctx, product, 5*time.Minute); err != nil {
		log.Printf("Failed to update product cache: %v", err)
	}

	return nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Удаляем из кеша
	if err := s.cache.Delete(ctx, id); err != nil {
		log.Printf("Failed to delete product from cache: %v", err)
	}

	return nil
}

// Метод для фонового обновления кеша
func (s *ProductService) StartCacheUpdater(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Product cache updater stopped")
			return
		case <-ticker.C:
			s.updateCache(ctx)
		}
	}
}

func (s *ProductService) updateCache(ctx context.Context) {
	log.Println("Updating product cache...")
	products, err := s.repo.GetAll(ctx)
	if err != nil {
		log.Printf("Failed to get products for cache update: %v", err)
		return
	}

	if err := s.cache.SetAllProducts(ctx, products, 10*time.Minute); err != nil {
		log.Printf("Failed to update products cache: %v", err)
		return
	}

	log.Printf("Product cache updated with %d products", len(products))
}
