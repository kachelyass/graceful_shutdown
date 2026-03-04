package database

import (
	"context"

	"example.com/m/v2/internal/models"
	"github.com/jmoiron/sqlx"
)

type TaskStore struct {
	db *sqlx.DB
}

func NewTaskStore(db *sqlx.DB) *TaskStore {
	return &TaskStore{db: db}
}

func (s *TaskStore) GetAll(ctx context.Context) ([]models.Task, error) {
	var tasks []models.Task

	query := `
SELECT * FROM tasks`
	err := s.db.SelectContext(ctx, &tasks, query)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
