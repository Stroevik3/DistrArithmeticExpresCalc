package memstore

import (
	"container/list"
	"sync"

	"github.com/Stroevik3/DistrArithmeticExpresCalc/internal/app/orchestratorserv/model"
	"github.com/Stroevik3/DistrArithmeticExpresCalc/internal/app/orchestratorserv/store"
	"github.com/google/uuid"
)

type Store struct {
	expressionRepository     *ExpressionRepository
	taskRepository           *TaskRepository
	StackReadyTaskRepository *StackReadyTaskRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) Expression() store.ExpressionRepository {
	if s.expressionRepository != nil {
		return s.expressionRepository
	}

	s.expressionRepository = &ExpressionRepository{
		store:       s,
		expressions: make(map[int]*model.Expression),
	}

	return s.expressionRepository
}

func (s *Store) Task() store.TaskRepository {
	if s.taskRepository != nil {
		return s.taskRepository
	}

	s.taskRepository = &TaskRepository{
		store: s,
		tasks: make(map[uuid.UUID]*model.Task),
	}

	return s.taskRepository
}

func (s *Store) StackReadyTask() store.StackReadyTaskRepository {
	if s.StackReadyTaskRepository != nil {
		return s.StackReadyTaskRepository
	}

	s.StackReadyTaskRepository = &StackReadyTaskRepository{
		mu:    sync.Mutex{},
		store: s,
		stack: list.New(),
	}

	return s.StackReadyTaskRepository
}
