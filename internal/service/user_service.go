package service

import (
	"context"
	"github.com/SaveljevRoman/go-layout-project/internal/models"
	"log"
	"time"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetAll(ctx context.Context) ([]*models.User, error)
	Create(ctx context.Context, user *models.User) (int64, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int64) error
}

type UserCache interface {
	GetByID(ctx context.Context, id int64) (*models.User, error)
	Set(ctx context.Context, user *models.User, expiration time.Duration) error
	Delete(ctx context.Context, id int64) error
	SetAllUsers(ctx context.Context, users []*models.User, expiration time.Duration) error
}

type UserService struct {
	repo  UserRepository
	cache UserCache
}

func NewUserService(repo UserRepository, cache UserCache) *UserService {
	return &UserService{
		repo:  repo,
		cache: cache,
	}
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*models.User, error) {
	// Сначала пытаемся получить из кеша
	user, err := s.cache.GetByID(ctx, id)
	if err != nil {
		log.Printf("Cache error: %v", err)
	}

	if user != nil {
		return user, nil
	}

	// Если в кеше нет, получаем из БД
	user, err = s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user != nil {
		// Кешируем результат на 5 минут
		if err := s.cache.Set(ctx, user, 5*time.Minute); err != nil {
			log.Printf("Failed to cache user: %v", err)
		}
	}

	return user, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) (int64, error) {
	id, err := s.repo.Create(ctx, user)
	if err != nil {
		return 0, err
	}

	// Обновить пользователя с ID
	user.ID = id
	if err := s.cache.Set(ctx, user, 5*time.Minute); err != nil {
		log.Printf("Failed to cache user: %v", err)
	}

	return id, nil
}

func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}

	// Обновляем кеш
	if err := s.cache.Set(ctx, user, 5*time.Minute); err != nil {
		log.Printf("Failed to update user cache: %v", err)
	}

	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Удаляем из кеша
	if err := s.cache.Delete(ctx, id); err != nil {
		log.Printf("Failed to delete user from cache: %v", err)
	}

	return nil
}

// StartCacheUpdater Метод для фонового обновления кеша
func (s *UserService) StartCacheUpdater(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Cache updater stopped")
			return
		case <-ticker.C:
			s.updateCache(ctx)
		}
	}
}

func (s *UserService) updateCache(ctx context.Context) {
	log.Println("Updating cache...")
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		log.Printf("Failed to get users for cache update: %v", err)
		return
	}

	if err := s.cache.SetAllUsers(ctx, users, 10*time.Minute); err != nil {
		log.Printf("Failed to update users cache: %v", err)
		return
	}

	log.Printf("Cache updated with %d users", len(users))
}
