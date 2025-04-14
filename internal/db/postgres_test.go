package pg_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pg "github.com/volchok96/todoapp/internal/db"
	"github.com/volchok96/todoapp/internal/domain"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&domain.Task{})
	require.NoError(t, err)
	return db
}

func newRepo(t *testing.T) domain.TaskRepository {
	db := setupTestDB(t)
	logger := zap.NewNop()
	return pg.NewTaskRepo(db, logger)
}

func TestTaskRepo_CreateAndGet(t *testing.T) {
	repo := newRepo(t)

	task := &domain.Task{Title: "Test", Date: "2025-04-14"}
	err := repo.Create(task)
	require.NoError(t, err)
	assert.NotZero(t, task.ID)

	fetched, err := repo.GetByID(task.ID)
	require.NoError(t, err)
	assert.Equal(t, task.Title, fetched.Title)
}

func TestTaskRepo_Update(t *testing.T) {
	repo := newRepo(t)

	task := &domain.Task{Title: "Initial", Date: "2025-04-14"}
	err := repo.Create(task)
	require.NoError(t, err)

	task.Title = "Updated"
	err = repo.Update(task)
	require.NoError(t, err)

	fetched, _ := repo.GetByID(task.ID)
	assert.Equal(t, "Updated", fetched.Title)
}

func TestTaskRepo_Update_NonExistent(t *testing.T) {
	repo := newRepo(t)

	err := repo.Update(&domain.Task{ID: 999, Title: "Missing"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "task not found")
}

func TestTaskRepo_Delete(t *testing.T) {
	repo := newRepo(t)

	task := &domain.Task{Title: "To be deleted", Date: "2025-04-14"}
	err := repo.Create(task)
	require.NoError(t, err)

	err = repo.Delete(task.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(task.ID)
	assert.Error(t, err)
}

func TestTaskRepo_Delete_NonExistent(t *testing.T) {
	repo := newRepo(t)
	err := repo.Delete(999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no task deleted")
}

func TestTaskRepo_List(t *testing.T) {
	repo := newRepo(t)

	require.NoError(t, repo.Create(&domain.Task{Title: "A", Date: "2025-04-14"}))
	require.NoError(t, repo.Create(&domain.Task{Title: "B", Date: "2025-04-14", Done: true}))

	tasks, err := repo.List(0, 10, nil)
	require.NoError(t, err)
	assert.Len(t, tasks, 2)

	done := true
	tasks, err = repo.List(0, 10, &done)
	require.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.True(t, tasks[0].Done)

	tasks, err = repo.List(0, 0, nil)
	require.NoError(t, err)
	assert.Len(t, tasks, 0)
}

func TestTaskRepo_ListByDate(t *testing.T) {
	repo := newRepo(t)

	require.NoError(t, repo.Create(&domain.Task{Title: "A", Date: "2025-04-14"}))
	require.NoError(t, repo.Create(&domain.Task{Title: "B", Date: "2025-04-14", Done: true}))
	require.NoError(t, repo.Create(&domain.Task{Title: "C", Date: "2025-04-15"}))

	tasks, err := repo.ListByDate("2025-04-14", nil)
	require.NoError(t, err)
	assert.Len(t, tasks, 2)

	done := true
	tasks, err = repo.ListByDate("2025-04-14", &done)
	require.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.Equal(t, "B", tasks[0].Title)
}

func TestTaskRepo_ListByDate_InvalidFormat(t *testing.T) {
	repo := newRepo(t)

	tasks, err := repo.ListByDate("not-a-date", nil)
	assert.NoError(t, err)
	assert.Len(t, tasks, 0)
}
