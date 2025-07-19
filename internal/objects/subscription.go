package objects

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Основная структура системы
type Subscription struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey"`     // уникальный идендификатор
	ServiceName string     `gorm:"not null"`                 // Название сервиса
	Price       int        `gorm:"not null;check:price > 0"` // Цена подписки
	UserID      uuid.UUID  `gorm:"type:uuid;not null"`       // уникальный id пользователя
	StartDate   time.Time  `gorm:"not null"`                 // Начало активации подписки
	EndDate     *time.Time // Окончание подписки
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
