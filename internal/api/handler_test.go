package api

import (
	"bytes"
	"context"
	"effective_mobile/internal/objects"
	"effective_mobile/pkg/logger_module"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSubscriptionService struct {
	mock.Mock
}

func (m *MockSubscriptionService) Create(ctx context.Context, sub *objects.Subscription) error {
	args := m.Called(ctx, sub)
	return args.Error(0)
}

func (m *MockSubscriptionService) GetByID(ctx context.Context, id uuid.UUID) (*objects.Subscription, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*objects.Subscription), args.Error(1)
}

func (m *MockSubscriptionService) Update(ctx context.Context, id uuid.UUID, fields map[string]interface{}) error {
	args := m.Called(ctx, id, fields)
	return args.Error(0)
}

func (m *MockSubscriptionService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSubscriptionService) Get_List(ctx context.Context, limit, offset int) ([]*objects.Subscription, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*objects.Subscription), args.Error(1)
}

func (m *MockSubscriptionService) GetTotalCost(ctx context.Context, userID uuid.UUID, serviceName string, start, end time.Time) (int, error) {
	args := m.Called(ctx, userID, serviceName, start, end)
	return args.Int(0), args.Error(1)
}

func TestCreateSubscription_Success(t *testing.T) {
	// Подготавливаем моки
	mockService := new(MockSubscriptionService)
	logger := logger_module.Get()

	handler := &SubscriptionHandler{
		service: mockService,
		logger:  logger,
	}

	test_req := &objects.Subscription{
		ServiceName: "Netflix",
		Price:       500,
		UserID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		StartDate:   time.Date(2025, time.November, 1, 0, 0, 0, 0, time.UTC),
	}

	mockService.On("Create", mock.Anything, test_req).Return(nil)

	test_body := `{
	"service_name": "Netflix",
	"price": 500,
	"user_id": "550e8400-e29b-41d4-a716-446655440000",
	"start_date": "11-2025"
	}`

	request_test := httptest.NewRequest("POST", "/api/subscriptions", bytes.NewBufferString(test_body))
	w := httptest.NewRecorder()

	handler.CreateSubscription(w, request_test)

	// Проверяем статус ответа
	assert.Equal(t, http.StatusCreated, w.Code)

	// Тут проверяем структуру ответа
	var sub_test objects.Subscription
	err := json.NewDecoder(w.Body).Decode(&sub_test)
	assert.NoError(t, err)

	assert.Equal(t, 500, sub_test.Price)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", sub_test.UserID.String())
	assert.Equal(t, "Netflix", sub_test.ServiceName)

	// Проверка что мок вызван со всеми аргументами
	mockService.AssertExpectations(t)
}

// Проверка невалидного json
func TestCreateSubscription_InvalidDate(t *testing.T) {
	mockService := new(MockSubscriptionService)
	logger := logger_module.New(nil, "[TEST]", 0)

	handler := &SubscriptionHandler{
		service: mockService,
		logger:  logger,
	}

	test_body := `{
		"service_name": "Netflix",
		"price": 500,
		"user_id": "550e8400-e29b-41d4-a716-446655440000",
		"start_date": "invalid-date"
	}`

	request_test := httptest.NewRequest("POST", "/api/subscriptions", bytes.NewBufferString(test_body))
	w := httptest.NewRecorder()

	handler.CreateSubscription(w, request_test)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid format start_data")
}

func TestGetSubscription(t *testing.T) {
	// Настройка моков
	mock_service := new(MockSubscriptionService)
	logger := logger_module.Get()

	handler := &SubscriptionHandler{
		service: mock_service,
		logger:  logger,
	}

	// Готовим тестовые данные
	testID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	test_Sub := &objects.Subscription{
		ID:          testID,
		ServiceName: "Netflix",
		Price:       599,
		UserID:      uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba"),
		StartDate:   time.Date(2025, time.July, 1, 0, 0, 0, 0, time.UTC),
	}
	// Настройка ожидания
	mock_service.On("GetByID", mock.Anything, testID).Return(test_Sub, nil)

	// Создаем тут тестовый запрос
	request_test := httptest.NewRequest("GET", "/api/subscriptions/"+testID.String(), nil)
	w := httptest.NewRecorder()

	// Добавляем параметр ID в запрос (для mux.Vars)
	request_test = mux.SetURLVars(request_test, map[string]string{"id": testID.String()})

	// Вызываем функцию обработчика
	handler.GetSubscription(w, request_test)

	// Проверка статуса кода
	assert.Equal(t, http.StatusOK, w.Code)
	var response objects.Subscription
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	// Тут проверяем структуру ответа
	assert.Equal(t, test_Sub.ID, response.ID)
	assert.Equal(t, test_Sub.ServiceName, response.ServiceName)
	assert.Equal(t, test_Sub.Price, response.Price)
	mock_service.AssertExpectations(t)
}

func TestGetSubscription_InvalidID(t *testing.T) {
	mockService := new(MockSubscriptionService)
	logger := logger_module.Get()

	handler := &SubscriptionHandler{
		service: mockService,
		logger:  logger,
	}

	// Добавим неккоректный uuid подписки
	invalidID := "invalid-uuid"

	request_test := httptest.NewRequest("GET", "/api/subscriptions/"+invalidID, nil)
	w := httptest.NewRecorder()

	// Добавляем параметр ID в запрос
	request_test = mux.SetURLVars(request_test, map[string]string{"id": invalidID})

	handler.GetSubscription(w, request_test)

	// Проверки
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid subscription id")
	mockService.AssertNotCalled(t, "GetByID")
}

func TestGetListSubscription_Success(t *testing.T) {
	mock_service := new(MockSubscriptionService)
	logger := logger_module.Get()

	handler := &SubscriptionHandler{
		service: mock_service,
		logger:  logger,
	}
	// Готовим тестовые данные
	test_list_subscription := []*objects.Subscription{
		{
			ID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			ServiceName: "Yandex",
			Price:       800,
			UserID:      uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba"),
			StartDate:   time.Date(2025, time.July, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
			ServiceName: "Yandex",
			Price:       800,
			UserID:      uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cbb"),
			StartDate:   time.Date(2025, time.July, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	// Настройка ожидания
	mock_service.On("Get_List", mock.Anything, 2, 0).Return(test_list_subscription, nil)

	// Создаем тестовый запрос
	request_test := httptest.NewRequest("GET", "/api/subscriptions?limit=2&offset=0", nil)
	w := httptest.NewRecorder()

	// Вызываем обработчик
	handler.GetListSubscription(w, request_test)

	// Проверки
	assert.Equal(t, http.StatusOK, w.Code)

	var response []objects.Subscription
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Len(t, response, 2)
	assert.Equal(t, test_list_subscription[0].ServiceName, response[0].ServiceName)
	assert.Equal(t, test_list_subscription[1].Price, response[1].Price)
	mock_service.AssertExpectations(t)
}

func TestGetListSubscription_InvalidParams(t *testing.T) {
	mockService := new(MockSubscriptionService)
	logger := logger_module.Get()

	handler := &SubscriptionHandler{
		service: mockService,
		logger:  logger,
	}

	testCases := []struct {
		name         string
		url          string
		expectLimit  int
		expectOffset int
	}{
		{"Negative limit", "/subscriptions?limit=-1&offset=0", -1, 0},
		{"Negative offset", "/subscriptions?limit=10&offset=-5", 10, -5},
		{"Non-numeric limit", "/subscriptions?limit=abc&offset=0", 0, 0},
		{"Non-numeric offset", "/subscriptions?limit=10&offset=xyz", 10, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Настраиваем ожидание
			mockService.On("Get_List", mock.Anything, tc.expectLimit, tc.expectOffset).
				Return([]*objects.Subscription{}, nil)

			request_test := httptest.NewRequest("GET", tc.url, nil)
			w := httptest.NewRecorder()

			handler.GetListSubscription(w, request_test)

			assert.Equal(t, http.StatusOK, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestUpdateSubscription_Success(t *testing.T) {
	mockService := new(MockSubscriptionService)
	logger := logger_module.Get()

	handler := &SubscriptionHandler{
		service: mockService,
		logger:  logger,
	}

	// Тестовые данные
	testID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	body_test := `{
        "service_name": "New Service Name",
        "price": 600,
        "end_date": "12-2025"
    }`

	fields_for_update := map[string]interface{}{
		"service_name": "New Service Name",
		"price":        600,
		"end_date":     time.Date(2025, time.December, 1, 0, 0, 0, 0, time.UTC),
	}
	// Создаем ожидаемый результат
	mockService.On("Update", mock.Anything, testID, fields_for_update).Return(nil)

	//  Создаем тестовый запрос
	request_test := httptest.NewRequest("PATCH", "/api/subscriptions/"+testID.String(), bytes.NewBufferString(body_test))

	// Добавляем параметр ID в запрос (для mux.Vars)
	request_test = mux.SetURLVars(request_test, map[string]string{"id": testID.String()})
	w := httptest.NewRecorder()

	handler.UpdateSubscription(w, request_test)

	// Проверка статуса кода
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверка стурктуры ответа
	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])

	// Проверка что мок вызван со всеми аргументами
	mockService.AssertExpectations(t)
}

func TestUpdateSubscription_InvalidID(t *testing.T) {
	mockService := new(MockSubscriptionService)
	logger := logger_module.Get()

	handler := &SubscriptionHandler{
		service: mockService,
		logger:  logger,
	}
	// Создаем тестовый запрос с невалидный id для проверки работы
	request_test := httptest.NewRequest("PATCH", "/api/subscriptions/invalid-id", bytes.NewBufferString(`{}`))
	request_test = mux.SetURLVars(request_test, map[string]string{"id": "invalid-id"})
	w := httptest.NewRecorder()

	handler.UpdateSubscription(w, request_test)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid subscription id")
	mockService.AssertNotCalled(t, "Update")
}

// Тест на проверку запроса без полей
func TestUpdateSubscription_NoFiled(t *testing.T) {
	mock_service := new(MockSubscriptionService)
	logger := logger_module.Get()

	handler := &SubscriptionHandler{
		service: mock_service,
		logger:  logger,
	}

	testID := uuid.New()
	request_test := httptest.NewRequest("PATCH", "/api/subscription/"+testID.String(), bytes.NewBufferString("{}"))
	request_test = mux.SetURLVars(request_test, map[string]string{"id": testID.String()})
	w := httptest.NewRecorder()

	handler.UpdateSubscription(w, request_test)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Так как у нас возращается json ответ, мы проверяем через map
	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "no fields for update", response["error"])
	mock_service.AssertNotCalled(t, "Update")
}

func TestDeleteSubscription_Success(t *testing.T) {
	mockService := new(MockSubscriptionService)
	logger := logger_module.Get()

	handler := &SubscriptionHandler{
		service: mockService,
		logger:  logger,
	}
	// Создаем тестовый id
	testID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	mockService.On("Delete", mock.Anything, testID).Return(nil)

	// Создаем тестовый запрос
	request_test := httptest.NewRequest("DELETE", "/subscriptions/"+testID.String(), nil)
	request_test = mux.SetURLVars(request_test, map[string]string{"id": testID.String()})
	w := httptest.NewRecorder()

	handler.DeleteSubscription(w, request_test)

	// Проверка
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Contains(t, w.Body.String(), "null")
	mockService.AssertExpectations(t)
}

// Проверка невалидного id
func TestDeleteSubscription_InvalidID(t *testing.T) {
	mockService := new(MockSubscriptionService)
	logger := logger_module.Get()

	handler := &SubscriptionHandler{
		service: mockService,
		logger:  logger,
	}
	// Создание тестового запроса на удаление с невалидным id
	request_test := httptest.NewRequest("DELETE", "/subscriptions/invalid-id", nil)
	request_test = mux.SetURLVars(request_test, map[string]string{"id": "invalid-id"})
	w := httptest.NewRecorder()

	handler.DeleteSubscription(w, request_test)

	// Проверка
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid subscription ID")
	mockService.AssertNotCalled(t, "Delete")
}
func TestDeleteSubscription_NotFound(t *testing.T) {
	mockService := new(MockSubscriptionService)
	logger := logger_module.Get()

	handler := &SubscriptionHandler{
		service: mockService,
		logger:  logger,
	}

	testID := uuid.New()

	// Настрайваем ожидание
	mockService.On("Delete", mock.Anything, testID).Return(errors.New("not found"))

	// Создание тестового запроса на удаление с невалидным id
	request_test := httptest.NewRequest("DELETE", "/subscriptions/"+testID.String(), nil)
	request_test = mux.SetURLVars(request_test, map[string]string{"id": testID.String()})
	w := httptest.NewRecorder()

	handler.DeleteSubscription(w, request_test)

	// Проверка
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "subscription not found")
	mockService.AssertExpectations(t)
}
