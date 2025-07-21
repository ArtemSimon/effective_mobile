package repository

import (
	"context"
	"effective_mobile/internal/objects"
	"time"

	"github.com/google/uuid"
)

// Реализуем интерфейс для работы с БД
type SubsctriptionRepository interface {
	Create(ctx context.Context, subscription *objects.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*objects.Subscription, error)
	Update(ctx context.Context, id uuid.UUID, fields map[string]interface{}) error
	Delete(ctx context.Context, id uuid.UUID) error
	Get_List(ctx context.Context, limit, offset int) ([]*objects.Subscription, error)
	GetTotalCost(
		ctx context.Context,
		userID uuid.UUID,
		service_name string,
		start_time, end_time time.Time,
	) (int, error)
}
