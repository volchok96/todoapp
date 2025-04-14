package usecase

import "github.com/volchok96/todoapp/internal/domain"

type taskUsecase struct {
	repo domain.TaskRepository
}

func NewTaskUsecase(repo domain.TaskRepository) domain.TaskUsecase {
	return &taskUsecase{repo: repo}
}

func (uc *taskUsecase) Create(task *domain.Task) error {
	return uc.repo.Create(task)
}

func (uc *taskUsecase) GetByID(id uint) (*domain.Task, error) {
	return uc.repo.GetByID(id)
}

func (uc *taskUsecase) Update(task *domain.Task) error {
	return uc.repo.Update(task)
}

func (uc *taskUsecase) Delete(id uint) error {
	return uc.repo.Delete(id)
}

func (uc *taskUsecase) List(offset, limit int, status *bool) ([]*domain.Task, error) {
	return uc.repo.List(offset, limit, status)
}

func (uc *taskUsecase) ListByDate(date string, status *bool) ([]*domain.Task, error) {
	return uc.repo.ListByDate(date, status)
}