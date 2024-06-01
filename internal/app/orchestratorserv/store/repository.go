package store

import (
	"github.com/Stroevik3/DistrArithmeticExpresCalc/internal/app/orchestratorserv/model"
	"github.com/google/uuid"
)

type ExpressionRepository interface {
	Add(*model.Expression) error
	GetList(bool) ([]*model.Expression, error)
	Find(int) (*model.Expression, error)
}

type TaskRepository interface {
	Add(*model.Task) error
	GetList() ([]*model.Task, error)
	Find(uuid.UUID) (*model.Task, error)
}

type StackReadyTaskRepository interface {
	Push(t *model.Task) error
	Pop() (*model.Task, error)
}
