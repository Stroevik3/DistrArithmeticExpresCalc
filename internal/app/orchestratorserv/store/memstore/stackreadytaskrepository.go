package memstore

import (
	"container/list"
	"errors"
	"sync"

	"github.com/Stroevik3/DistrArithmeticExpresCalc/internal/app/orchestratorserv/model"
	"github.com/Stroevik3/DistrArithmeticExpresCalc/internal/app/orchestratorserv/store"
)

type StackReadyTaskRepository struct {
	mu    sync.Mutex
	store *Store
	stack *list.List
}

func (r *StackReadyTaskRepository) Push(t *model.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if t.Oper == "/" && t.ArgTaskTwo == nil && t.ArgTwo == 0 {
		t.Status = TASK_STATUS_ERORR
		return errors.New("you can't divide by zero")
	}
	t.Status = TASK_STATUS_READY
	r.stack.PushBack(t)
	return nil
}

func (r *StackReadyTaskRepository) Pop() (*model.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.stack.Len() == 0 {
		return nil, store.ErrStackIsEmpty
	}
	e := r.stack.Front()
	t := e.Value.(*model.Task)
	r.stack.Remove(e)
	t.Status = TASK_STATUS_PROC
	return t, nil
}
