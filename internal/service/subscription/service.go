package subscription

import (
	"Price/internal/domain/product"
	sub "Price/internal/domain/subscription"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Service interface {
	Create(ctx context.Context, subscription sub.Subscription) (sub.Subscription, error)
	ListAll(ctx context.Context, userID int64) ([]product.Product, error)
	GetSubscribedUsers(ctx context.Context, prod int64) ([]int64, error)
	GetAll(ctx context.Context) ([]product.Product, error)
	UpdatePrice(ctx context.Context, userID int64, prodID int64, price float64) error
	Delete(ctx context.Context, userID int64, prodID int64) error
	GetSub(ctx context.Context, userID int64, productID int64) (sub.Subscription, error)
}

type SubscriptionService struct {
	repo sub.Repository

	r *redis.Client
}

func NewService(repo sub.Repository, r *redis.Client) *SubscriptionService {
	return &SubscriptionService{
		repo: repo,
		r:    r,
	}
}

func (s *SubscriptionService) getSubCacheKey(userID, prodID int64) string {
	return fmt.Sprintf("sub:%d:%d", userID, prodID)
}

func (s *SubscriptionService) Create(ctx context.Context, subscription sub.Subscription) (sub.Subscription, error) {
	if subscription.UserID <= 0 {
		return sub.Subscription{}, errors.New("неверный ID пользователя")
	}
	if subscription.ProductID <= 0 {
		return sub.Subscription{}, errors.New("неверный ID продукта")
	}
	if subscription.TargetPrice <= 0 {
		return sub.Subscription{}, errors.New("неверная ожидаемая цена продукта")
	}
	cacheKey := s.getSubCacheKey(subscription.UserID, subscription.ProductID)
	s.r.Del(ctx, cacheKey)
	return s.repo.Create(ctx, subscription)
}

func (s *SubscriptionService) ListAll(ctx context.Context, userID int64) ([]product.Product, error) {

	return s.repo.GetProduct(ctx, userID)
}

func (s *SubscriptionService) GetSubscribedUsers(ctx context.Context, prod int64) ([]int64, error) {
	return s.repo.GetSubscribedUsers(ctx, prod)
}

func (s *SubscriptionService) GetAll(ctx context.Context) ([]product.Product, error) {
	return s.repo.GetAll(ctx)
}

func (s *SubscriptionService) UpdatePrice(ctx context.Context, userID int64, prodID int64, price float64) error {
	err := s.repo.UpdatePrice(ctx, userID, prodID, price)
	cacheKey := s.getSubCacheKey(userID, prodID)
	if err == nil {
		s.r.Del(ctx, cacheKey)
	}
	return err
}

func (s *SubscriptionService) Delete(ctx context.Context, userID int64, prodID int64) error {
	if userID <= 0 {
		return errors.New("ID пользователя не может быть нулевым")
	}
	if prodID <= 0 {
		return errors.New("ID продукта не может быть нулевым")
	}

	err := s.repo.Delete(ctx, userID, prodID)
	if err != nil {
		return err
	}
	cacheKey := s.getSubCacheKey(userID, prodID)
	s.r.Del(ctx, cacheKey)
	return err
}

func (s *SubscriptionService) GetSub(ctx context.Context, userID int64, productID int64) (sub.Subscription, error) {
	if userID <= 0 {
		return sub.Subscription{}, errors.New("ID пользователя не может быть нулевым")
	}
	if productID <= 0 {
		return sub.Subscription{}, errors.New("ID продукта не может быть нулевым")
	}

	cacheKey := s.getSubCacheKey(userID, productID)

	cachedData, err := s.r.Get(ctx, cacheKey).Bytes()
	if err == redis.Nil {
		log.Printf("[Redis] Промах кэша для ключа %s. Идем в БД.", cacheKey)

		// Достаем из Postgres
		subscription, dbErr := s.repo.GetSubscription(ctx, userID, productID)
		if dbErr != nil {
			return sub.Subscription{}, dbErr
		}

		bytesToCache, marshalErr := json.Marshal(subscription)
		if marshalErr == nil {
			s.r.Set(ctx, cacheKey, bytesToCache, 30*time.Minute)
		}
		return subscription, nil
	}
	if err != nil {
		log.Printf("ошибка redis: %v", err)
		return s.repo.GetSubscription(ctx, userID, productID)
	}
	log.Printf("удалось достать данные из redis:%s", cacheKey)
	var subscription sub.Subscription
	if unmarshalErr := json.Unmarshal(cachedData, &subscription); unmarshalErr != nil {
		return sub.Subscription{}, unmarshalErr
	}

	return subscription, nil
}
func (s *SubscriptionService) GetUsersForPriceDrop(ctx context.Context, productID int64, currentPrice float64) ([]int64, error) {
	return s.GetUsersForPriceDrop(ctx, productID, currentPrice)
}
