package memstore

import (
	"errors"
	"sync"

	"github.com/Stroevik3/DistrArithmeticExpresCalc/internal/app/orchestratorserv/model"
	"github.com/Stroevik3/DistrArithmeticExpresCalc/internal/app/orchestratorserv/store"
	"github.com/google/uuid"
)

const (
	TASK_STATUS_CREATE     string = "CREATE"           // Создан
	TASK_STATUS_READY      string = "READY_TO_PROCESS" // Готов к обработке
	TASK_STATUS_PROC       string = "PROCESSING"       // Обрабатывается
	TASK_STATUS_ERORR      string = "ERORR"            // Обрабатывается
	TASK_STATUS_RES_DETERM string = "RES_DETERMINED"   // Результат определен
)

type TaskRepository struct {
	mu    sync.Mutex
	store *Store
	tasks map[uuid.UUID]*model.Task
}

func (r *TaskRepository) Add(t *model.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	var err error
	if t.Oper == "/" && t.ArgTaskTwo == nil && t.ArgTwo == 0 {
		t.Status = TASK_STATUS_ERORR
		return errors.New("you can't divide by zero")
	}
	t.Id, err = uuid.NewUUID()
	if err != nil {
		return err
	}
	t.Status = TASK_STATUS_CREATE
	t.Result = 0
	r.tasks[t.Id] = t
	return nil
}

func (r *TaskRepository) GetList() ([]*model.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	newTasks := make([]*model.Task, 0, len(r.tasks))
	for _, val := range r.tasks {
		newTasks = append(newTasks, &model.Task{
			Id:         val.Id,
			Exp:        val.Exp,
			ArgOne:     val.ArgOne,
			ArgTaskOne: val.ArgTaskOne,
			ArgTwo:     val.ArgTwo,
			ArgTaskTwo: val.ArgTaskTwo,
			Oper:       val.Oper,
			OperTime:   val.OperTime,
			Prior:      val.Prior,
			Status:     val.Status,
			Result:     val.Result,
			TaskNext:   val.TaskNext,
		})
	}
	return newTasks, nil
}

func (r *TaskRepository) Find(id uuid.UUID) (*model.Task, error) {
	t, ok := r.tasks[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}
	return t, nil
}
