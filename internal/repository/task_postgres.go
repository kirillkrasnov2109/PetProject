package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"petProject/internal/domain"
)

type taskRepo struct {
	db *DB
}

func NewTaskRepository(db *DB) domain.TaskRepository {
	return &taskRepo{db: db}
}

var errTaskNotFound = errors.New("task not found") // Пользовательская ошибка

func (r *taskRepo) Create(ctx context.Context, req domain.TaskRequest) (int, error) {
	payload, err := marshalPayload(req.Payload)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal payload: %w", err)
	}

	const query = `INSERT INTO tasks (webhook, send_date, payload) VALUES ($1, $2, $3) RETURNING id`

	var id int
	err = r.db.pool.QueryRow(ctx, query, req.Webhook, req.SendDate, payload).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert task: %w", err)
	}

	return id, nil
}

func (r *taskRepo) Read(ctx context.Context, id int) (domain.Task, error) {
	const query = `SELECT id, webhook, send_date, payload FROM tasks WHERE id=$1`

	var t domain.Task
	var payload []byte

	err := r.db.pool.QueryRow(ctx, query, id).Scan(&t.Id, &t.Webhook, &t.SendDate, &payload)
	if err != nil {
		return domain.Task{}, fmt.Errorf("task not found: %w", err)
	}

	if err := unmarshalPayload(payload, &t.Payload); err != nil {
		return domain.Task{}, fmt.Errorf("invalid payload data: %w", err)
	}

	return t, nil
}

func (r *taskRepo) Update(ctx context.Context, t domain.Task) error {
	payload, err := marshalPayload(t.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	const query = `UPDATE tasks SET webhook=$1, send_date=$2, payload=$3 WHERE id=$4`

	result, err := r.db.pool.Exec(ctx, query, t.Webhook, t.SendDate, payload, t.Id)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errTaskNotFound
	}

	return nil
}

func (r *taskRepo) Delete(ctx context.Context, id int) error {
	const query = `DELETE FROM tasks WHERE id=$1`

	result, err := r.db.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errTaskNotFound
	}

	return nil
}

func (r *taskRepo) List(ctx context.Context) ([]domain.Task, error) {
	const query = `SELECT id, webhook, send_date, payload FROM tasks`

	rows, err := r.db.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var t domain.Task
		var payload []byte

		if err := rows.Scan(&t.Id, &t.Webhook, &t.SendDate, &payload); err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
		}
		if err := unmarshalPayload(payload, &t.Payload); err != nil {
			return nil, fmt.Errorf("payload unmarshalling failed: %w", err)
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

//MARK: - Payload helpers

func marshalPayload(payload interface{}) ([]byte, error) {
	return json.Marshal(payload)
}

func unmarshalPayload(data []byte, dest interface{}) error {
	return json.Unmarshal(data, dest)
}
