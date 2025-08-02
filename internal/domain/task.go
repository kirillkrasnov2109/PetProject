package domain

import (
	"context"
	"time"
)

type Task struct { // структура, которая включает все поля
	Id int `json:"id"`
	TaskRequest
}

// POST /task

type TaskRequest struct { // структура запроса
	Webhook  string    `json:"webhook"`
	SendDate time.Time `json:"send_date"`
	Payload  any       `json:"payload"`
}

// GET /task/:id
// Response Task

// PUT /task/:id
// Request TaskRequest

// DELETE /task/:id

// TODO:
// GET /task

type TaskRepository interface { // интерфейс для работы с БД
	Create(ctx context.Context, task TaskRequest) (int, error)
	Read(ctx context.Context, id int) (Task, error)
	Update(ctx context.Context, task Task) error
	Delete(ctx context.Context, id int) error

	List(ctx context.Context) ([]Task, error)
}

type TaskService interface { // интерфейс для бизнес логики
	Create(ctx context.Context, task TaskRequest) (int, error)
	Read(ctx context.Context, id int) (Task, error)
	Update(ctx context.Context, id int, task TaskRequest) error
	Delete(ctx context.Context, id int) error

	List(ctx context.Context) ([]Task, error)
}
