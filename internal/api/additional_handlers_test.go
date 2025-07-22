package api

import (
	"effective_mobile/pkg/logger_module"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetTotalCost_Success(t *testing.T) {
	mockService := new(MockSubscriptionService)
	logger := logger_module.Get()

	handler := &SubscriptionHandler{
		service: mockService,
		logger:  logger,
	}

	// Тестовые данные
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	serviceName := "Netflix"
	startDate := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, time.December, 1, 0, 0, 0, 0, time.UTC)
	result_total := 12000

	// Настраиваем ожидание
	mockService.On("GetTotalCost", mock.Anything, userID, serviceName, startDate, endDate).
		Return(result_total, nil)

	// Создаем тестовый запрос
	request_test := httptest.NewRequest("GET",
		"/api/subscriptions/total?user_id="+userID.String()+
			"&service_name="+serviceName+
			"&start=01-2025"+
			"&end=12-2025", nil)
	w := httptest.NewRecorder()

	handler.GetTotalCost(w, request_test)

	// Проверка
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]int
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, result_total, response["total"])
	mockService.AssertExpectations(t)
}
func TestGetTotalCost_InvalidStartDate(t *testing.T) {
	mockService := new(MockSubscriptionService)
	logger := logger_module.Get()

	handler := &SubscriptionHandler{
		service: mockService,
		logger:  logger,
	}

	// Генерим userID
	userID := uuid.New()

	// Создаем тестовый запрос с невалидной startdata
	request_test := httptest.NewRequest("GET",
		"/api/subscriptions/total?user_id="+userID.String()+
			"&service_name=Netflix"+
			"&start=invalid"+
			"&end=12-2025", nil)
	w := httptest.NewRecorder()

	handler.GetTotalCost(w, request_test)

	// Проверка
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid start date format")
	mockService.AssertNotCalled(t, "GetTotalCost")
}

func TestGetTotalCost_InvalidEndDate(t *testing.T) {
	mockService := new(MockSubscriptionService)
	logger := logger_module.Get()

	handler := &SubscriptionHandler{
		service: mockService,
		logger:  logger,
	}
	// Генерим userID
	userID := uuid.New()
	// Создаем тестовый запрос с невалидной enddata
	request_test := httptest.NewRequest("GET",
		"/api/subscriptions/total?user_id="+userID.String()+
			"&service_name=Netflix"+
			"&start=01-2025"+
			"&end=invalid", nil)
	w := httptest.NewRecorder()

	handler.GetTotalCost(w, request_test)

	// Проверка
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid end date format")
	mockService.AssertNotCalled(t, "GetTotalCost")
}

func TestGetTotalCost_ServiceError(t *testing.T) {
	mockService := new(MockSubscriptionService)
	logger := logger_module.Get()

	handler := &SubscriptionHandler{
		service: mockService,
		logger:  logger,
	}

	// Генерим userID
	userID := uuid.New()
	startDate := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, time.December, 1, 0, 0, 0, 0, time.UTC)

	// Настройка ожидания
	mockService.On("GetTotalCost", mock.Anything, userID, "Netflix", startDate, endDate).
		Return(0, errors.New("calculation error"))

	// Создаем тестовый запрос
	request_test := httptest.NewRequest("GET",
		"/api/subscriptions/total?user_id="+userID.String()+
			"&service_name=Netflix"+
			"&start=01-2025"+
			"&end=12-2025", nil)
	w := httptest.NewRecorder()

	handler.GetTotalCost(w, request_test)

	// Проверка
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "internal server error")
	mockService.AssertExpectations(t)
}
