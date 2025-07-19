package objects

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Реализуем интерфейс для работы с БД
type SubsctriptionRepository interface {
	Create(ctx context.Context, s *Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*Subscription, error)
	Update(ctx context.Context, s *Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	Get_List(ctx context.Context, limit, offset int) ([]*Subscription, error)
	GetTotalCost(
		ctx context.Context,
		userID uuid.UUID,
		service_name string,
		start_time, end_time time.Time,
	) (int, error)
	AutoMigrate(ctx context.Context) error
}
