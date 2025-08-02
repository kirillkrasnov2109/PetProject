package service

import (
	"context"
	"petProject/internal/domain"
)

type taskService struct {
	repo domain.TaskRepository
}

func NewTaskService(repo domain.TaskRepository) domain.TaskService {
	return &taskService{repo: repo}
}

func (s *taskService) Create(ctx context.Context, req domain.TaskRequest) (int, error) {
	return s.repo.Create(ctx, req)
}

func (s *taskService) Read(ctx context.Context, id int) (domain.Task, error) {
	return s.repo.Read(ctx, id)
}

func (s *taskService) Update(ctx context.Context, id int, req domain.TaskRequest) error {
	return s.repo.Update(ctx, domain.Task{
		Id:          id,
		TaskRequest: req,
	})
}

func (s *taskService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *taskService) List(ctx context.Context) ([]domain.Task, error) {
	return s.repo.List(ctx)
}
