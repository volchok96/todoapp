package http

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/volchok96/todoapp/internal/domain"
	"go.uber.org/zap"
)

func RegisterRoutes(app *fiber.App, uc domain.TaskUsecase, logger *zap.Logger) {
	r := app.Group("/tasks")

	r.Post("/", CreateTaskHandler(uc, logger))
	r.Get("/:id", GetTaskHandler(uc, logger))
	r.Put("/:id", UpdateTaskHandler(uc, logger))
	r.Delete("/:id", DeleteTaskHandler(uc, logger))
	r.Get("/", ListTasksHandler(uc, logger))
	r.Get("/by-date/:date", ListTasksByDateHandler(uc, logger))
}

// CreateTaskHandler godoc
// @Summary Create a task
// @Description Create a new task
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body domain.Task true "Task to create"
// @Success 201 {object} domain.Task
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /tasks [post]
func CreateTaskHandler(uc domain.TaskUsecase, logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		task := new(domain.Task)
		if err := c.BodyParser(task); err != nil {
			logger.Error("Failed to parse JSON", zap.Error(err))
			return c.Status(400).JSON(fiber.Map{"error": "cannot parse JSON"})
		}
		if err := uc.Create(task); err != nil {
			logger.Error("Failed to create task", zap.Error(err))
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		logger.Info("Task created", zap.String("title", task.Title))
		return c.Status(201).JSON(task)
	}
}

// GetTaskHandler godoc
// @Summary Get task by ID
// @Description Get a single task by its ID
// @Tags tasks
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} domain.Task
// @Failure 404 {object} domain.ErrorResponse
// @Router /tasks/{id} [get]
func GetTaskHandler(uc domain.TaskUsecase, logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		task, err := uc.GetByID(uint(id))
		if err != nil {
			logger.Error("Task not found", zap.Uint("id", uint(id)), zap.Error(err))
			return c.Status(404).JSON(fiber.Map{"error": "task not found"})
		}
		return c.JSON(task)
	}
}

// UpdateTaskHandler godoc
// @Summary Update task
// @Description Update an existing task
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param task body domain.Task true "Updated task"
// @Success 200 {object} domain.Task
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /tasks/{id} [put]
func UpdateTaskHandler(uc domain.TaskUsecase, logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		task := new(domain.Task)
		if err := c.BodyParser(task); err != nil {
			logger.Error("Failed to parse JSON", zap.Error(err))
			return c.Status(400).JSON(fiber.Map{"error": "cannot parse JSON"})
		}
		task.ID = uint(id)
		if err := uc.Update(task); err != nil {
			logger.Error("Failed to update task", zap.Error(err))
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(task)
	}
}

// DeleteTaskHandler godoc
// @Summary Delete task
// @Description Delete a task by ID
// @Tags tasks
// @Param id path int true "Task ID"
// @Success 204
// @Failure 500 {object} domain.ErrorResponse
// @Router /tasks/{id} [delete]
func DeleteTaskHandler(uc domain.TaskUsecase, logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		if err := uc.Delete(uint(id)); err != nil {
			logger.Error("Failed to delete task", zap.Uint("id", uint(id)), zap.Error(err))
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(204)
	}
}

// ListTasksHandler godoc
// @Summary List tasks
// @Description List tasks with pagination and optional status filter
// @Tags tasks
// @Produce json
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Param done query bool false "Status filter"
// @Success 200 {array} domain.Task
// @Failure 500 {object} domain.ErrorResponse
// @Router /tasks [get]
func ListTasksHandler(uc domain.TaskUsecase, logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		offset, _ := strconv.Atoi(c.Query("offset", "0"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))
		statusParam := c.Query("done")
		var status *bool
		if statusParam != "" {
			val := statusParam == "true"
			status = &val
		}
		tasks, err := uc.List(offset, limit, status)
		if err != nil {
			logger.Error("Failed to fetch task list", zap.Error(err))
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(tasks)
	}
}

// ListTasksByDateHandler godoc
// @Summary List tasks by date
// @Description List tasks for a specific date and optional status
// @Tags tasks
// @Produce json
// @Param date path string true "Date in YYYY-MM-DD"
// @Param done query bool false "Status filter"
// @Success 200 {array} domain.Task
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /tasks/by-date/{date} [get]
func ListTasksByDateHandler(uc domain.TaskUsecase, logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		dateStr := c.Params("date")

		// Валидация формата даты
		if _, err := time.Parse("2006-01-02", dateStr); err != nil {
			logger.Error("Invalid date format",
				zap.String("received", dateStr),
				zap.Error(err))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid date format. Expected YYYY-MM-DD",
			})
		}

		// Обработка параметра done
		var status *bool
		if doneParam := c.Query("done"); doneParam != "" {
			val, err := strconv.ParseBool(doneParam)
			if err != nil {
				logger.Error("Invalid done parameter",
					zap.String("received", doneParam),
					zap.Error(err))
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid done parameter. Expected true/false",
				})
			}
			status = &val
		}

		tasks, err := uc.ListByDate(dateStr, status)
		if err != nil {
			logger.Error("Failed to fetch tasks",
				zap.String("date", dateStr),
				zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch tasks",
			})
		}

		return c.JSON(tasks)
	}
}
