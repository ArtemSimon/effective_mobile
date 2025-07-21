package repository

import (
	"context"
	"effective_mobile/internal/objects"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormRepo struct {
	db *gorm.DB
}

// Принимает готовое подключение *gorm.DB
// Возвращает реализацию репозитория
func NewGormRepo(db *gorm.DB) SubsctriptionRepository {
	return &GormRepo{db: db}
}

// Сохраняет подписку по id в БД
func (gr *GormRepo) Create(ctx context.Context, subscription *objects.Subscription) error {
	new_subscription := gr.db.WithContext(ctx).Create(subscription) // Добавляет контекст к запросу .WithContext (позволяет отменить операцию)
	if new_subscription.Error != nil {
		return fmt.Errorf("database error: %w", new_subscription.Error)
	}
	return nil
}

// Удаляет подписку по id
func (gr *GormRepo) Delete(ctx context.Context, id uuid.UUID) error {
	subscription_del := gr.db.WithContext(ctx).Delete(&objects.Subscription{}, "id = ?", id)
	if subscription_del.Error != nil {
		return fmt.Errorf("database error: %w", subscription_del.Error)
	}
	return nil
}

// Получаем список подписок с пагинацией
// SELECT * FROM subscriptions LIMIT {limit} OFFSET {offset};
func (gr *GormRepo) Get_List(ctx context.Context, limit, offset int) ([]*objects.Subscription, error) {

	var subscriptions []*objects.Subscription
	subscription_list := gr.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Find(&subscriptions)
	if subscription_list.Error != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", subscription_list.Error)
	}
	return subscriptions, nil
}

// Получаем подписку по id
// SELECT * FROM subscriptions WHERE id = '...' LIMIT 1;
func (gr *GormRepo) GetByID(ctx context.Context, id uuid.UUID) (*objects.Subscription, error) {
	var subscription objects.Subscription
	subscription_by_id := gr.db.WithContext(ctx).First(&subscription, "id = ?", id)
	if subscription_by_id.Error == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("subscription %w with id: %s", gorm.ErrRecordNotFound, id)
	}
	if subscription_by_id.Error != nil {
		return nil, fmt.Errorf("failed to get subcription: %w", subscription_by_id.Error)
	}
	return &subscription, nil
}

// Обновляем подписку по конкретным полям
// UPDATE subscriptions
// SET field1 = value1, field2 = value2
// WHERE id = 'ваш-uuid';
func (gr *GormRepo) Update(ctx context.Context, id uuid.UUID, fields map[string]interface{}) error {
	update_subscription := gr.db.WithContext(ctx).
		Model(&objects.Subscription{}).
		Where("id = ?", id).
		Updates(fields)

	if errors.Is(update_subscription.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("subscription %w with id: %s", update_subscription.Error, id)
	}
	if update_subscription.Error != nil {
		return fmt.Errorf("failed to update subscription: %w", update_subscription.Error)
	}

	return nil
}

// Реализуем кастомную функцию миграции для гибкости и так же инкапсулируем реализацию GORM
// func (gr *GormRepo) AutoMigrate(ctx context.Context) error {
// 	return gr.db.WithContext(ctx).AutoMigrate(&objects.Subscription{})
// }

// Реализуем функцию для просморта суммы подписок за определенный период
// SELECT COALESCE(SUM(price), 0)
// FROM subscriptions
// WHERE user_id = '...'
//
//	AND service_name = '...'
//	AND start_date >= '2023-01-01'
//	AND (end_date <= '2023-12-31' OR end_date IS NULL)
func (gr *GormRepo) GetTotalCost(
	ctx context.Context,
	userID uuid.UUID,
	serviceName string,
	start, end time.Time,
) (int, error) {
	var total int64

	query := gr.db.WithContext(ctx).
		Model(&objects.Subscription{}).
		Select("COALESCE(SUM(price), 0)").
		Where("(start_date <= ? AND (end_date >= ? OR end_date IS NULL))", end, start)

	if userID != uuid.Nil {
		query = query.Where("user_id = ?", userID)
	}
	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}

	if err := query.Scan(&total).Error; err != nil {
		return 0, fmt.Errorf("failed to calculate total cost: %w", err)
	}

	return int(total), nil
}
