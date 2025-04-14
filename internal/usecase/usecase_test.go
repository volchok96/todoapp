package usecase_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/volchok96/todoapp/internal/domain"
	"github.com/volchok96/todoapp/internal/usecase"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) Create(task *domain.Task) error {
	args := m.Called(task)
	return args.Error(0)
}
func (m *mockRepo) GetByID(id uint) (*domain.Task, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Task), args.Error(1)
}
func (m *mockRepo) Update(task *domain.Task) error {
	args := m.Called(task)
	return args.Error(0)
}
func (m *mockRepo) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *mockRepo) List(offset, limit int, status *bool) ([]*domain.Task, error) {
	args := m.Called(offset, limit, status)
	return args.Get(0).([]*domain.Task), args.Error(1)
}
func (m *mockRepo) ListByDate(date string, status *bool) ([]*domain.Task, error) {
	args := m.Called(date, status)
	return args.Get(0).([]*domain.Task), args.Error(1)
}

func TestTaskUsecase_Create(t *testing.T) {
	repo := new(mockRepo)
	uc := usecase.NewTaskUsecase(repo)

	task := &domain.Task{Title: "Test"}
	repo.On("Create", task).Return(nil)

	err := uc.Create(task)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestTaskUsecase_GetByID_NotFound(t *testing.T) {
	repo := new(mockRepo)
	uc := usecase.NewTaskUsecase(repo)

	repo.On("GetByID", uint(99)).Return(&domain.Task{}, errors.New("not found"))

	_, err := uc.GetByID(99)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}
