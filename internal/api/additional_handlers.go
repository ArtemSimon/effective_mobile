package api

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// GetTotalCost возвращает суммарную стоимость подписок
// @Summary Подсчет стоимости
// @Description Подсчитываем суммарную стоимость всех подписок за выбранный период
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id query string false "ID пользователя (UUID) для фильтрации" example("550e8400-e29b-41d4-a716-446655440000")
// @Param service_name query string false "Название сервиса для фильтрации" example("Netflix")
// @Param start query string true "Начало периода (формат MM-YYYY)" example("01-2025")
// @Param end query string true "Конец периода (формат MM-YYYY)" example("10-2025")
// @Success 200 {object} TotalCostResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/subscriptions/total [get]
func (handler *SubscriptionHandler) GetTotalCost(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("GetTotalCost handler called", "method", r.Method, "path", r.URL.Path)
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	handler.logger.Debug("Getting params from query")
	params := r.URL.Query()
	userID, _ := uuid.Parse(params.Get("user_id"))

	serviceName := params.Get("service_name")

	handler.logger.Debug("Parsing param start date")
	start, err := time.Parse("01-2006", params.Get("start"))
	if err != nil {
		handler.logger.Error("Failed parse start invalid date format",
			"error", err.Error(),
			"status_code", http.StatusBadRequest)
		sendError(w, http.StatusBadRequest, "invalid start date format")
		return
	}
	handler.logger.Debug("Parsing param end date")
	end, err := time.Parse("01-2006", params.Get("end"))
	if err != nil {
		handler.logger.Error("Failed parse end invalid date format",
			"error", err.Error(),
			"status_code", http.StatusBadRequest)
		sendError(w, http.StatusBadRequest, "invalid end date format")
		return
	}
	handler.logger.Debug("Calling service to get total cost")

	total, err := handler.service.GetTotalCost(ctx, userID, serviceName, start, end)
	if err != nil {
		handler.logger.Error("Failed get total cost for subscrioptions",
			"error", err.Error(),
			"status_code", http.StatusBadRequest)
		sendError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	handler.logger.Info("Successfully get total cost")
	renderJSON(w, http.StatusOK, map[string]int{"total": total})
}

// TotalCostResponse структура ответа для суммы подписок
type TotalCostResponse struct {
	Total int `json:"total"`
}
