package service

import (
	"context"
	"errors"
	"github.com/SaveljevRoman/go-layout-project/internal/models"
	"log"
	"time"
)

type PurchaseRepository interface {
	GetByID(ctx context.Context, id int64) (*models.Purchase, error)
	GetByUserID(ctx context.Context, userID int64) ([]*models.Purchase, error)
	Create(ctx context.Context, purchase *models.Purchase) (int64, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
	GetAll(ctx context.Context) ([]*models.Purchase, error)
}

type PurchaseCache interface {
	GetByID(ctx context.Context, id int64) (*models.Purchase, error)
	Set(ctx context.Context, purchase *models.Purchase, expiration time.Duration) error
	Delete(ctx context.Context, id int64) error
	SetUserPurchases(ctx context.Context, userID int64, purchases []*models.Purchase, expiration time.Duration) error
	GetUserPurchases(ctx context.Context, userID int64) ([]*models.Purchase, error)
}

type PurchaseService struct {
	repo           PurchaseRepository
	cache          PurchaseCache
	userService    *UserService
	productService *ProductService
}

func NewPurchaseService(repo PurchaseRepository, cache PurchaseCache, userService *UserService, productService *ProductService) *PurchaseService {
	return &PurchaseService{
		repo:           repo,
		cache:          cache,
		userService:    userService,
		productService: productService,
	}
}

func (s *PurchaseService) CreatePurchase(ctx context.Context, request *models.PurchaseRequest) (*models.Purchase, error) {
	// Проверяем существование пользователя
	user, err := s.userService.GetUser(ctx, request.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("пользователь не найден")
	}

	// Проверяем существование товара
	product, err := s.productService.GetProduct(ctx, request.ProductID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("товар не найден")
	}

	// Проверяем, что количество товара больше нуля
	if request.Quantity <= 0 {
		return nil, errors.New("количество товара должно быть больше нуля")
	}

	// Рассчитываем общую стоимость
	totalPrice := product.Price * float64(request.Quantity)

	// Создаем покупку
	purchase := &models.Purchase{
		UserID:     request.UserID,
		ProductID:  request.ProductID,
		Quantity:   request.Quantity,
		TotalPrice: totalPrice,
		Status:     "pending",
	}

	// Сохраняем в БД
	id, err := s.repo.Create(ctx, purchase)
	if err != nil {
		return nil, err
	}

	// Обновляем покупку с ID
	purchase.ID = id
	purchase.CreatedAt = time.Now()
	purchase.UpdatedAt = time.Now()

	// Кешируем результат
	if err := s.cache.Set(ctx, purchase, 5*time.Minute); err != nil {
		log.Printf("Failed to cache purchase: %v", err)
	}

	// Инвалидируем кеш пользовательских покупок
	s.cache.Delete(ctx, id)

	return purchase, nil
}

func (s *PurchaseService) GetPurchase(ctx context.Context, id int64) (*models.Purchase, error) {
	// Сначала пытаемся получить из кеша
	purchase, err := s.cache.GetByID(ctx, id)
	if err != nil {
		log.Printf("Cache error: %v", err)
	}

	if purchase != nil {
		return purchase, nil
	}

	// Если в кеше нет, получаем из БД
	purchase, err = s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if purchase != nil {
		// Кешируем результат на 5 минут
		if err := s.cache.Set(ctx, purchase, 5*time.Minute); err != nil {
			log.Printf("Failed to cache purchase: %v", err)
		}
	}

	return purchase, nil
}

func (s *PurchaseService) GetUserPurchases(ctx context.Context, userID int64) ([]*models.Purchase, error) {
	// Сначала пытаемся получить из кеша
	purchases, err := s.cache.GetUserPurchases(ctx, userID)
	if err != nil {
		log.Printf("Cache error: %v", err)
	}

	if purchases != nil && len(purchases) > 0 {
		return purchases, nil
	}

	// Если в кеше нет, получаем из БД
	purchases, err = s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Кешируем результат на 5 минут
	if len(purchases) > 0 {
		if err := s.cache.SetUserPurchases(ctx, userID, purchases, 5*time.Minute); err != nil {
			log.Printf("Failed to cache user purchases: %v", err)
		}
	}

	return purchases, nil
}

func (s *PurchaseService) UpdatePurchaseStatus(ctx context.Context, id int64, status string) error {
	// Проверяем допустимость статуса
	if status != "pending" && status != "completed" && status != "cancelled" {
		return errors.New("недопустимый статус покупки")
	}

	// Обновляем статус в БД
	if err := s.repo.UpdateStatus(ctx, id, status); err != nil {
		return err
	}

	// Обновляем кеш
	purchase, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if purchase != nil {
		if err := s.cache.Set(ctx, purchase, 5*time.Minute); err != nil {
			log.Printf("Failed to update purchase cache: %v", err)
		}
	}

	return nil
}

func (s *PurchaseService) GetAllPurchases(ctx context.Context) ([]*models.Purchase, error) {
	return s.repo.GetAll(ctx)
}

// StartCacheUpdater Метод для фонового обновления кеша покупок
func (s *PurchaseService) StartCacheUpdater(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Purchase cache updater stopped")
			return
		case <-ticker.C:
			s.updateCache(ctx)
		}
	}
}

func (s *PurchaseService) updateCache(ctx context.Context) {
	log.Println("Updating purchase cache...")
	// В данной реализации мы просто логируем событие обновления кеша
	// Полная реализация может включать получение всех покупок и их кеширование
	// Но это может быть ресурсоемко, поэтому лучше кешировать по мере запросов
}
