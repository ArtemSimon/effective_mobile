package api

import (
	"context"
	"effective_mobile/internal/objects"
	"effective_mobile/internal/service"
	"effective_mobile/pkg/logger_module"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type SubscriptionHandler struct {
	service service.SubscriptionServiceI
	logger  *logger_module.Logger
}

func NewSubciptionHandler(service service.SubscriptionServiceI, logger *logger_module.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{service: service, logger: logger}
}

// Структура для ответа при обновлении
type UpdateResponce struct {
	Status int `json:"status"`
}

// Структура для ошибок в API
type ErrorResponse struct {
	Error string `json:"error"`
}

func renderJSON(w http.ResponseWriter, code int, object interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(object)
}

// Для отправки ошибок
func sendError(w http.ResponseWriter, code int, message string) {
	renderJSON(w, code, ErrorResponse{Error: message})
}

// Данная ручка создает новую подписку
// @Summary Создать подписку
// @Description Создать новую запись о подписке пользователя
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param sub body objects.SubscriptionCreateRequest true "Данные подписки"
// @Success 201 {object} objects.Subscription
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/subscriptions [post]
func (handler *SubscriptionHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 7*time.Second)
	defer cancel()

	var req_sub objects.SubscriptionCreateRequest

	handler.logger.Debug("Decode request body")
	err := json.NewDecoder(r.Body).Decode(&req_sub)
	if err != nil {
		handler.logger.Error("failed to request body", "error", err.Error(), "status_code", http.StatusBadRequest)
		sendError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	handler.logger.Debug("Request body decoded successfully",
		"service_name", req_sub.ServiceName,
		"user_id", req_sub.UserID)

	handler.logger.Debug("Started parse start date", "start_date", req_sub.StartDate)
	start_Date, err := time.Parse("01-2006", req_sub.StartDate)
	if err != nil {
		handler.logger.Error("Invalid start date format",
			"error", err.Error(),
			"start_date", req_sub.StartDate,
			"status_code", http.StatusBadRequest)
		sendError(w, http.StatusBadRequest, "invalid format start_data")
		return
	}

	handler.logger.Debug("Started parse user ID", "user_id", req_sub.UserID)
	user_ID, err := uuid.Parse(req_sub.UserID)
	if err != nil {
		handler.logger.Error("Invalid user ID format",
			"error", err.Error(),
			"user_id", req_sub.UserID,
			"status_code", http.StatusBadRequest)
		sendError(w, http.StatusBadRequest, "invalid user_id format")
		return
	}

	sub := &objects.Subscription{
		ServiceName: req_sub.ServiceName,
		Price:       req_sub.Price,
		UserID:      user_ID,
		StartDate:   start_Date,
	}

	if req_sub.EndDate != nil {
		handler.logger.Debug("Start parse end date", "end_date", *req_sub.EndDate)
		end_Date, err := time.Parse("01-2006", *req_sub.EndDate)
		if err != nil {
			handler.logger.Error("Invalid end date format",
				"error", err.Error(),
				"end date", *req_sub.EndDate,
				"status_code", http.StatusBadRequest)
			sendError(w, http.StatusBadRequest, "invalid end_date format")
			return
		}
		sub.EndDate = &end_Date
	}

	handler.logger.Info("Creating subscription",
		"service_name", sub.ServiceName,
		"user_id", sub.UserID,
		"start_date", sub.StartDate.Format("01-2006"),
		"price", sub.Price)

	handler.logger.Debug("Calling service to create subscription")

	if err := handler.service.Create(ctx, sub); err != nil {
		handler.logger.Error("Failed to create subscription",
			"error", err.Error(),
			"status_code", http.StatusInternalServerError)
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	handler.logger.Info("Subscription created successfully",
		"service_name", sub.ServiceName)
	renderJSON(w, http.StatusCreated, sub)
}

// Данная ручка возвращает подписку по ID
// @Summary Получить подписку
// @Description Получаем подписку по id
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки" format(uuid) example("550e8400-e29b-41d4-a716-446655440000")
// @Success 200 {object} objects.Subscription
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/subscriptions/{id} [get]
func (handler *SubscriptionHandler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("GetSubscription handler called", "method", r.Method,
		"path", r.URL.Path)

	ctx, cancel := context.WithTimeout(r.Context(), 7*time.Second)
	defer cancel()

	handler.logger.Debug("Start parse subscription id")
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		handler.logger.Error("Invalid subscription ID format",
			"error", err.Error(),
			"subscription_id", id,
			"status_code", http.StatusBadRequest)
		sendError(w, http.StatusBadRequest, "invalid subscription id")
		return
	}
	handler.logger.Debug("Calling service to get subscription",
		"subscription_id", id)

	sub, err := handler.service.GetByID(ctx, id)
	if err != nil {
		handler.logger.Error("Failed to retrieve subscription",
			"error", err.Error(),
			"subscription_id", id,
			"status_code", http.StatusNotFound)
		sendError(w, http.StatusNotFound, "subscription not found")
		return
	}

	handler.logger.Info("Subscription successfully get",
		"subscription_id", id,
		"service_name", sub.ServiceName,
		"user_id", sub.UserID)
	renderJSON(w, http.StatusOK, sub)
}

// Данная ручка возвращает список подписок с пагинацией
// @Summary Получаем подписки
// @Description Получаем все подписки которые есть
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param limit query integer false "Лимит записей (по умолчанию 10)"
// @Param offset query integer false "Смещение (по умолчанию 0)"
// @Success 200 {array} objects.Subscription
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/subscriptions [get]
func (handler *SubscriptionHandler) GetListSubscription(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("GetListSubscription handler called", "method", r.Method,
		"path", r.URL.Path)
	ctx, cancel := context.WithTimeout(r.Context(), 7*time.Second)
	defer cancel()

	handler.logger.Debug("Parse params limit,offset in query")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	handler.logger.Debug("Calling service to get list subscription",
		"params", "limit", limit, "offset", offset)

	subscriptions, err := handler.service.Get_List(ctx, limit, offset)
	if err != nil {
		handler.logger.Error("Failed to get all subscriptions",
			"error", err.Error(),
			"status_code", http.StatusInternalServerError)
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	handler.logger.Info("Successfully get list subscrition", "params", "limit", limit, "offset", offset)
	renderJSON(w, http.StatusOK, subscriptions)
}

// Данная ручка обновляет подписку
// @Summary Обновляем подписку
// @Description Обновляем подписку по указанному полю
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки формата UUID" format(uuid) example("550e8400-e29b-41d4-a716-446655440000")
// @Param request body objects.SubscriptionUpdateRequest true "Поля для обновления"
// @Success 200 {object} UpdateResponce
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/subscriptions/{id} [patch]
func (handler *SubscriptionHandler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("UpdateSubscription handler called", "method", r.Method,
		"path", r.URL.Path)
	ctx, cancel := context.WithTimeout(r.Context(), 7*time.Second)
	defer cancel()

	// Получаем ID из пути
	handler.logger.Info("Started parse id in path")
	variable := mux.Vars(r)
	id, err := uuid.Parse(variable["id"])
	if err != nil {
		handler.logger.Error("Invalid format for subscription id",
			"error", err.Error(),
			"status_code", http.StatusBadRequest)
		sendError(w, http.StatusBadRequest, "invalid subscription id")
		return
	}

	var updateStruct objects.SubscriptionUpdateRequest
	handler.logger.Debug("Decode request body")
	if err := json.NewDecoder(r.Body).Decode(&updateStruct); err != nil {
		handler.logger.Error("failed to request body", "error", err, "status_code", http.StatusBadRequest)
		sendError(w, http.StatusBadRequest, "invalid request fields")
		return
	}

	handler.logger.Debug("Request body decoded successfully")

	// Проверяем, что есть хотя бы одно поле для обновления
	if updateStruct.ServiceName == nil && updateStruct.Price == nil && updateStruct.EndDate == nil {
		handler.logger.Error("No fields for update", "error", err, http.StatusBadRequest)
		sendError(w, http.StatusBadRequest, "no fields for update")
		return
	}
	handler.logger.Debug("Start create map for fields from request body")
	// Преобразуем в map для GORM
	fields := make(map[string]interface{})
	if updateStruct.ServiceName != nil {
		handler.logger.Debug("Create field service name")
		fields["service_name"] = *updateStruct.ServiceName
	}
	if updateStruct.Price != nil {
		handler.logger.Debug("Create field price")
		fields["price"] = *updateStruct.Price
	}
	if updateStruct.EndDate != nil {
		// Преобразуем строку даты в time.Time
		handler.logger.Debug("Parse and create field endDate")
		if endDate, err := time.Parse("01-2006", *updateStruct.EndDate); err == nil {
			fields["end_date"] = endDate
		}
	}
	handler.logger.Debug("Calling service to delete subscription by id",
		"subscription_id", id, "fields", fields)

	if err := handler.service.Update(ctx, id, fields); err != nil {
		handler.logger.Error("Failed update subscription by fields",
			"error", err.Error(),
			"status_code", http.StatusNotFound)
		sendError(w, http.StatusNotFound, "subscription not found")
		return
	}
	handler.logger.Info("Successfully update subscription")
	renderJSON(w, http.StatusOK, map[string]string{"status": "success"})

}

// Данная ручка удаляет подписку
// @Summary Удаление подписки
// @Description Удаляем подписку по указанному id
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки в формате UUID" format(uuid) example("550e8400-e29b-41d4-a716-446655440000")
// @Success 204 "Подписка успешно удалена"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/subscriptions/{id} [delete]
func (handler *SubscriptionHandler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("DeleteSubscription handler called", "method", r.Method,
		"path", r.URL.Path)
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	handler.logger.Debug("Start parse id")
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		handler.logger.Error("Invalid subscription id format",
			"error", err.Error(),
			"subscription_id", id,
			"status_code", http.StatusBadRequest)
		sendError(w, http.StatusBadRequest, "invalid subscription ID")
		return
	}
	handler.logger.Debug("Calling service to delete subscription by id",
		"subscription_id", id)

	if err := handler.service.Delete(ctx, id); err != nil {
		handler.logger.Error("Failed delete subscription by id",
			"error", err.Error(),
			"status_code", http.StatusNotFound)
		sendError(w, http.StatusNotFound, "subscription not found")
		return
	}
	handler.logger.Info("Successfully delete subscription by id", "id", id)
	renderJSON(w, http.StatusNoContent, nil)
}
