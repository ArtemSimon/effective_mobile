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
	// logger.Println("GetTotalCost handler called")
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	params := r.URL.Query()
	userID, _ := uuid.Parse(params.Get("user_id"))

	serviceName := params.Get("service_name")

	start, err := time.Parse("01-2006", params.Get("start"))
	if err != nil {
		sendError(w, http.StatusBadRequest, "invalid start date format")
		return
	}

	end, err := time.Parse("01-2006", params.Get("end"))
	if err != nil {
		sendError(w, http.StatusBadRequest, "invalid end date format")
		return
	}

	total, err := handler.service.GetTotalCost(ctx, userID, serviceName, start, end)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	renderJSON(w, http.StatusOK, map[string]int{"total": total})
}

// TotalCostResponse структура ответа для суммы подписок
type TotalCostResponse struct {
	Total int `json:"total"`
}
