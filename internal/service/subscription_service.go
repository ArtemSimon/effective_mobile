package service

import (
	"context"
	"effective_mobile/internal/objects"
	"effective_mobile/internal/repository"
	"effective_mobile/pkg/logger_module"
	"time"

	"github.com/google/uuid"
)

// Интерфейс для сервисного слоя
type SubscriptionServiceI interface {
	Create(ctx context.Context, sub *objects.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*objects.Subscription, error)
	Update(ctx context.Context, id uuid.UUID, fields map[string]interface{}) error
	Delete(ctx context.Context, id uuid.UUID) error
	Get_List(ctx context.Context, limit, offset int) ([]*objects.Subscription, error)
	GetTotalCost(ctx context.Context, userID uuid.UUID, serviceName string, start, end time.Time) (int, error)
}

type SubscriptionService struct {
	rep    repository.SubsctriptionRepository // принимает обьект удовлетворяющий указанному interface, тут мы используем GormRepo
	logger *logger_module.Logger
}

func NewSubciptionService(rep repository.SubsctriptionRepository, logger *logger_module.Logger) SubscriptionServiceI {
	return &SubscriptionService{rep: rep, logger: logger}
}

func (subservice *SubscriptionService) Create(ctx context.Context, sub *objects.Subscription) error {
	if sub.Price <= 0 {
		subservice.logger.Fatal("price must be positive")
	}
	if sub.ServiceName == "" {
		subservice.logger.Fatal("service name is required")
	}
	subservice.logger.Debug("Calling db layer for create subscription")
	return subservice.rep.Create(ctx, sub)
}
func (subservice *SubscriptionService) GetByID(ctx context.Context, id uuid.UUID) (*objects.Subscription, error) {
	subservice.logger.Debug("Calling db layer for get subscription by id")
	return subservice.rep.GetByID(ctx, id)
}

func (subservice *SubscriptionService) Update(ctx context.Context, id uuid.UUID, fields map[string]interface{}) error {
	if price, ok := fields["price"].(int); ok && price <= 0 {
		subservice.logger.Fatal("price must be positive")
	}
	subservice.logger.Debug("Calling db layer for update subscription by fields")
	return subservice.rep.Update(ctx, id, fields)
}

func (subservice *SubscriptionService) Delete(ctx context.Context, id uuid.UUID) error {
	subservice.logger.Debug("Calling db layer for delete subscription by id")
	return subservice.rep.Delete(ctx, id)
}

func (subservice *SubscriptionService) Get_List(ctx context.Context, limit, offset int) ([]*objects.Subscription, error) {
	// Устанавливаем дефолтные значения
	subservice.logger.Info("Install default value for limit,offset")

	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	subservice.logger.Debug("Calling db layer for get all subscriptions")
	return subservice.rep.Get_List(ctx, limit, offset)
}

func (subservice *SubscriptionService) GetTotalCost(ctx context.Context, userID uuid.UUID, serviceName string, start, end time.Time) (int, error) {
	subservice.logger.Debug("Calling db layer for get total cost subscriptions")
	return subservice.rep.GetTotalCost(ctx, userID, serviceName, start, end)
}
