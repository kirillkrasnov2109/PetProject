package syncer

import (
	"context"
	"log"
	"petProject/internal/domain"
	"time"
)

type Syncer struct {
	service domain.TaskService
	sender  domain.Sender
}

func NewSyncer(service domain.TaskService, sender domain.Sender) *Syncer {
	return &Syncer{service: service, sender: sender}
}

func (s *Syncer) Run(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.runOnce(ctx)
		case <-ctx.Done():
			log.Println("syncer stopped")
			return
		}
	}
}

func (s *Syncer) runOnce(ctx context.Context) {
	tasks, err := s.service.List(ctx)
	if err != nil {
		log.Printf("failed to list tasks: %v", err)
		return
	}

	now := time.Now()
	for _, task := range tasks {
		if task.SendDate.Before(now) || task.SendDate.Equal(now) {
			go func(task domain.Task) {
				if err := s.sender.Send(task); err != nil {
					log.Printf("send failed: %v", err)
					return
				}
				_ = s.service.Delete(ctx, task.Id)
			}(task)
		}
	}
}
