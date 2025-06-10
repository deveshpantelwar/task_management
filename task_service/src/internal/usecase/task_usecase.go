package usecase

import (
	"context"
	"errors"
	"task_service/src/internal/core/session"
	"task_service/src/internal/core/task"
)

type TaskUsecase struct {
	repo task.Repository
}

func NewTaskUsecase(repo task.Repository) *TaskUsecase {
	return &TaskUsecase{repo: repo}
}

func (u *TaskUsecase) CreateTask(ctx context.Context, title, description, priority string, assignedTo *int64) (*session.Task, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}
	if priority == "" {
		priority = "medium"
	}
	t := &session.Task{
		Title:       title,
		Description: description,
		Status:      "pending",
		Priority:    priority,
		AssignedTo:  assignedTo,
	}
	err := u.repo.CreateTask(ctx, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (u *TaskUsecase) UpdateTask(ctx context.Context, taskID int, updatedTask *session.Task) error {
	return u.repo.UpdateTask(ctx, taskID, updatedTask)
}

func (u *TaskUsecase) ListTasks(ctx context.Context, assignedTo *int, status *string) ([]*session.Task, error) {
	return u.repo.ListTasks(ctx, assignedTo, status)
}

func (s *TaskUsecase) MarkTaskCompleted(ctx context.Context, taskID int64) error {
	return s.repo.MarkTaskCompleted(ctx, taskID)
}
