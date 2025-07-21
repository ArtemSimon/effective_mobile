package objects

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SubscriptionUpdateRequest определяет поля для обновления подписки

type SubscriptionUpdateRequest struct {
	ServiceName *string `json:"service_name,omitempty" example:"Netflix"`
	Price       *int    `json:"price,omitempty" example:"599"`
	EndDate     *string `json:"end_date,omitempty" example:"03-2025"`
}

// Отдельная структура для создания подписки
type SubscriptionCreateRequest struct {
	ServiceName string  `json:"service_name" example:"Netflix" binding:"required"`
	Price       int     `json:"price" example:"599" binding:"required"`
	UserID      string  `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" binding:"required"`
	StartDate   string  `json:"start_date" example:"09-2025" binding:"required"`
	EndDate     *string `json:"end_date" example:"03-2025"`
}

// Основная структура системы
type Subscription struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey"  json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`   // уникальный идендификатор
	ServiceName string     `gorm:"not null" json:"service_name" example:"Netflix"`                                   // Название сервиса
	Price       int        `gorm:"not null;check:price > 0" json:"price" example:"599"`                              // Цена подписки
	UserID      uuid.UUID  `gorm:"type:uuid;not null" json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"` // уникальный id пользователя
	StartDate   time.Time  `gorm:"not null" json:"start_date" swaggertype:"string" example:"09-2025"`                // Начало активации подписки
	EndDate     *time.Time `json:"end_date,omitempty" swaggertype:"string" example:"03-2025"`                        // Окончание подписки
}

func (s *Subscription) IsActive(time_subscription time.Time) bool {
	if s.EndDate == nil {
		return time_subscription.After(s.StartDate) || time_subscription.Equal(s.StartDate)
	}
	return time_subscription.After(s.StartDate) && time_subscription.Before(*s.EndDate)
}

// Хук перед созданием для генерации id если нету
func (s *Subscription) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
