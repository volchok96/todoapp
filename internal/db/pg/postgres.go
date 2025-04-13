package pg

import (
	"time"

	"github.com/volchok96/todoapp/internal/domain"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type taskRepo struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewTaskRepo(db *gorm.DB, logger *zap.Logger) domain.TaskRepository {
	return &taskRepo{db, logger}
}

func (r *taskRepo) Create(task *domain.Task) error {
	r.logger.Info("Creating task", zap.String("title", task.Title))
	err := r.db.Create(task).Error
	if err != nil {
		r.logger.Error("Failed to create task", zap.Error(err))
	}
	return err
}

func (r *taskRepo) GetByID(id uint) (*domain.Task, error) {
	r.logger.Info("Fetching task by ID", zap.Uint("id", id))
	var task domain.Task
	if err := r.db.First(&task, id).Error; err != nil {
		r.logger.Error("Failed to fetch task", zap.Error(err))
		return nil, err
	}
	return &task, nil
}

func (r *taskRepo) Update(task *domain.Task) error {
	r.logger.Info("Updating task", zap.Uint("id", task.ID), zap.String("title", task.Title))
	err := r.db.Save(task).Error
	if err != nil {
		r.logger.Error("Failed to update task", zap.Error(err))
	}
	return err
}

func (r *taskRepo) Delete(id uint) error {
	r.logger.Info("Deleting task", zap.Uint("id", id))
	err := r.db.Delete(&domain.Task{}, id).Error
	if err != nil {
		r.logger.Error("Failed to delete task", zap.Error(err))
	}
	return err
}

func (r *taskRepo) List(offset, limit int, status *bool) ([]*domain.Task, error) {
	var tasks []*domain.Task
	query := r.db.Offset(offset).Limit(limit)
	if status != nil {
		query = query.Where("done = ?", *status)
	}
	if err := query.Find(&tasks).Error; err != nil {
		r.logger.Error("Failed to fetch task list", zap.Error(err))
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepo) ListByDate(date string, status *bool) ([]*domain.Task, error) {
	var tasks []*domain.Task
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		r.logger.Error("Failed to parse date", zap.String("date", date), zap.Error(err))
		return nil, err
	}
	query := r.db.Where("date = ?", parsedDate)
	if status != nil {
		query = query.Where("done = ?", *status)
	}
	if err := query.Find(&tasks).Error; err != nil {
		r.logger.Error("Failed to fetch tasks by date", zap.String("date", date), zap.Error(err))
		return nil, err
	}
	return tasks, nil
}
