package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/volchok96/todoapp/internal/domain"
	httpDelivery "github.com/volchok96/todoapp/internal/http"
	"go.uber.org/zap"
)

type mockUsecase struct{}

func (m *mockUsecase) Create(task *domain.Task) error                     { return nil }
func (m *mockUsecase) GetByID(id uint) (*domain.Task, error)             { return &domain.Task{ID: id, Title: "Test", Description: "Test Desc"}, nil }
func (m *mockUsecase) Update(task *domain.Task) error                    { return nil }
func (m *mockUsecase) Delete(id uint) error                              { return nil }
func (m *mockUsecase) List(offset, limit int, status *bool) ([]*domain.Task, error) {
	return []*domain.Task{{ID: 1, Title: "Test"}}, nil
}
func (m *mockUsecase) ListByDate(date string, status *bool) ([]*domain.Task, error) {
	return []*domain.Task{{ID: 1, Title: "Test Date"}}, nil
}

func setupTestApp() *fiber.App {
	app := fiber.New()
	logger, _ := zap.NewDevelopment()
	uc := &mockUsecase{}
	httpDelivery.RegisterRoutes(app, uc, logger)
	return app
}

func TestGetTaskHandler(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest(http.MethodGet, "/tasks/42", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCreateTaskHandler(t *testing.T) {
	app := setupTestApp()

	task := domain.Task{
		Title:       "New Task",
		Description: "Test desc",
		Date:        "2024-04-14",
	}
	body, _ := json.Marshal(task)

	req := httptest.NewRequest(http.MethodPost, "/tasks/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

