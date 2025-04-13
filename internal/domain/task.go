package domain

import (
	"encoding/json"
	"time"
)

// swagger:response swaggerResponse
type ErrorResponse struct {
	// in:body
	Body struct {
		Error string `json:"error"`
	}
}

// Task represents a todo item
// swagger:model
type Task struct {
	// The id of the task
	// example: 1
	ID uint `json:"id" gorm:"primaryKey"`

	// The title of the task
	// example: Buy groceries
	Title string `json:"title"`

	// Detailed description of the task
	// example: Milk, eggs, bread
	Description string `json:"description"`

	// Due date in YYYY-MM-DD format
	// example: 2023-05-20
	Date string `json:"date" gorm:"type:date"`

	// Completion status
	// example: false
	Done bool `json:"done"`
}

type DateOnly struct {
	time.Time
}

const layout = "2006-01-02"

type TaskRepository interface {
	Create(task *Task) error
	GetByID(id uint) (*Task, error)
	Update(task *Task) error
	Delete(id uint) error
	List(offset, limit int, status *bool) ([]*Task, error)
	ListByDate(date string, status *bool) ([]*Task, error)
}

type TaskUsecase interface {
	Create(task *Task) error
	GetByID(id uint) (*Task, error)
	Update(task *Task) error
	Delete(id uint) error
	List(offset, limit int, status *bool) ([]*Task, error)
	ListByDate(date string, status *bool) ([]*Task, error)
}

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	str := string(b)
	str = str[1 : len(str)-1]
	t, err := time.Parse(layout, str)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Format(layout))
}
