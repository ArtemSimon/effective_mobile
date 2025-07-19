package repository

import (
	"context"
	"effective_mobile/internal/objects"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormRepo struct {
	db *gorm.DB
}

// Принимает готовое подключение *gorm.DB
// Возвращает реализацию репозитория
func NewGormRepo(db *gorm.DB) *GormRepo {
	return &GormRepo{db: db}
}

// Сохраняет подписку по id в БД
func (gr *GormRepo) Create(ctx context.Context, subscription *objects.Subscription) error {
	return gr.db.WithContext(ctx).Create(subscription).Error
}

// Удаляет подписку по id
func (gr *GormRepo) Delete(ctx context.Context, id uuid.UUID) error {
	result := gr.db.WithContext(ctx).Delete(&objects.Subscription{}, "id = ?", id)
	return result.Error
}

// Получаем список подписок с пагинацией
func (gr *GormRepo) Get_List(ctx context.Context, limit, offset int) ([]*objects.Subscription, error) {

	var subscriptions []*objects.Subscription
	result := gr.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Find(&subscriptions)
	return subscriptions, result.Error
}

// Получаем подписку по id
func (gr *GormRepo) GetByID(ctx context.Context, id uuid.UUID) (*objects.Subscription, error) {
	var subscription objects.Subscription
	subscription_by_id := gr.db.WithContext(ctx).First(&subscription, "id = ?", id)
	if subscription_by_id.Error == gorm.ErrRecordNotFound {
		return nil, gorm.ErrRecordNotFound
	}
	return &subscription, subscription_by_id.Error
}

// Обновляем подписку по конкретным полям
func (gr *GormRepo) Update(ctx context.Context, id uuid.UUID, fields map[string]interface{}) error {
	update_subscription := gr.db.WithContext(ctx).
		Model(&objects.Subscription{}).
		Where("id = ?", id).
		Updates(fields)
	return update_subscription.Error
}

// Реализуем кастомную функцию миграции для гибкости и так же инкапсулируем реализацию GORM
func (gr *GormRepo) AutoMigrate(ctx context.Context) error {
	return gr.db.WithContext(ctx).AutoMigrate(&objects.Subscription{})
}

// Реализуем функцию для просморта суммы подписок за определенный период
func (r *GormRepo) GetTotalCost(
	ctx context.Context,
	userID uuid.UUID,
	serviceName string,
	start, end time.Time,
) (int, error) {
	var total int64

	query := r.db.WithContext(ctx).
		Model(&objects.Subscription{}).
		Select("COALESCE(SUM(price), 0)")

	if userID != uuid.Nil {
		query = query.Where("user_id = ?", userID)
	}
	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}

	err := query.
		Where("start_date >= ?", start).
		Where("end_date <= ? OR end_date IS NULL", end).
		Scan(&total).Error

	return int(total), err
}
