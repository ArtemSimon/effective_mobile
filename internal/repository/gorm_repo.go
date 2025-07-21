package repository

import (
	"context"
	"effective_mobile/internal/objects"
	"effective_mobile/pkg/logger_module"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormRepo struct {
	db     *gorm.DB
	logger *logger_module.Logger
}

// Принимает готовое подключение *gorm.DB
// Возвращает реализацию репозитория
func NewGormRepo(db *gorm.DB, logger *logger_module.Logger) SubsctriptionRepository {
	return &GormRepo{db: db, logger: logger}
}

// Сохраняет подписку по id в БД
func (gr *GormRepo) Create(ctx context.Context, subscription *objects.Subscription) error {
	gr.logger.Info("Starting ORM request create subscription in db")
	new_subscription := gr.db.WithContext(ctx).Create(subscription) // Добавляет контекст к запросу .WithContext (позволяет отменить операцию)
	if new_subscription.Error != nil {
		gr.logger.Fatal("Database error", "error", new_subscription.Error)
	}
	return nil
}

// Удаляет подписку по id
func (gr *GormRepo) Delete(ctx context.Context, id uuid.UUID) error {
	gr.logger.Info("Starting ORM request delete subscription in db")
	subscription_del := gr.db.WithContext(ctx).Delete(&objects.Subscription{}, "id = ?", id)
	if subscription_del.Error != nil {
		gr.logger.Fatal("Database error", "error", subscription_del.Error)

	}
	return nil
}

// Получаем список подписок с пагинацией
// SELECT * FROM subscriptions LIMIT {limit} OFFSET {offset};
func (gr *GormRepo) Get_List(ctx context.Context, limit, offset int) ([]*objects.Subscription, error) {
	gr.logger.Info("Starting ORM request get list subscription in db")
	var subscriptions []*objects.Subscription
	subscription_list := gr.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Find(&subscriptions)
	if subscription_list.Error != nil {
		gr.logger.Fatal("Failed to get subscriptions", "error", subscription_list.Error)
	}
	gr.logger.Info("Successfully request in db to get list subscriptions")
	return subscriptions, nil
}

// Получаем подписку по id
// SELECT * FROM subscriptions WHERE id = '...' LIMIT 1;
func (gr *GormRepo) GetByID(ctx context.Context, id uuid.UUID) (*objects.Subscription, error) {
	gr.logger.Info("Starting ORM request get by id subscription in db")
	var subscription objects.Subscription
	subscription_by_id := gr.db.WithContext(ctx).First(&subscription, "id = ?", id)
	if subscription_by_id.Error == gorm.ErrRecordNotFound {
		gr.logger.Fatal("Failed subdcription not found", "error", gorm.ErrRecordNotFound, "id", id)
	}
	if subscription_by_id.Error != nil {
		gr.logger.Fatal("Failed to get subscription", "error", subscription_by_id.Error)
	}

	gr.logger.Info("Successfully request in db to get by id subscription")

	return &subscription, nil
}

// Обновляем подписку по конкретным полям
// UPDATE subscriptions
// SET field1 = value1, field2 = value2
// WHERE id = 'ваш-uuid';
func (gr *GormRepo) Update(ctx context.Context, id uuid.UUID, fields map[string]interface{}) error {
	gr.logger.Info("Starting ORM request update subscription in db")
	update_subscription := gr.db.WithContext(ctx).
		Model(&objects.Subscription{}).
		Where("id = ?", id).
		Updates(fields)

	if errors.Is(update_subscription.Error, gorm.ErrRecordNotFound) {
		gr.logger.Fatal("Failed update subscription", "error", update_subscription.Error, "id", id)
	}
	if update_subscription.Error != nil {
		gr.logger.Fatal("Failed to update subscription", "error", update_subscription.Error)
	}

	gr.logger.Info("Successfully request in db to update subscription")
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
	gr.logger.Info("Starting ORM request get total cost subscriptions in db")

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
		gr.logger.Fatal("Failed to calculate total cost", "error", err)
	}

	gr.logger.Info("Successfully request in db to get total cost subscriptions")

	return int(total), nil
}
