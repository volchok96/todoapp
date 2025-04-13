package docs

import "github.com/volchok96/todoapp/internal/domain"

// Generic response
// swagger:response swaggerResponse
type swaggerResponse struct {
	// in:body
	Body struct {
		Error string `json:"error"`
	}
}

// Task model
// swagger:response taskResponse
type taskResponse struct {
	// in:body
	Body domain.Task
}

// List of tasks
// swagger:response tasksResponse
type tasksResponse struct {
	// in:body
	Body []domain.Task
}