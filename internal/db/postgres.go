package pg

import (
	"fmt"

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

	var existing domain.Task
	if err := r.db.First(&existing, task.ID).Error; err != nil {
		r.logger.Warn("Task not found for update", zap.Uint("id", task.ID), zap.Error(err))
		return fmt.Errorf("task not found: %w", err)
	}

	if err := r.db.Save(task).Error; err != nil {
		r.logger.Error("Failed to update task", zap.Error(err))
		return err
	}
	return nil
}

func (r *taskRepo) Delete(id uint) error {
	r.logger.Info("Deleting task", zap.Uint("id", id))
	res := r.db.Delete(&domain.Task{}, id)
	if res.Error != nil {
		r.logger.Error("Failed to delete task", zap.Error(res.Error))
		return res.Error
	}
	if res.RowsAffected == 0 {
		r.logger.Warn("No rows affected during delete", zap.Uint("id", id))
		return fmt.Errorf("no task deleted")
	}
	return nil
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

func (r *taskRepo) ListByDate(dateStr string, status *bool) ([]*domain.Task, error) {
	var tasks []*domain.Task

	query := r.db.Where("date = ?", dateStr)
	if status != nil {
		query = query.Where("done = ?", *status)
	}

	if err := query.Find(&tasks).Error; err != nil {
		r.logger.Error("Database query failed",
			zap.String("date", dateStr),
			zap.Any("status", status),
			zap.Error(err))
		return nil, fmt.Errorf("database error: %w", err)
	}

	return tasks, nil
}
