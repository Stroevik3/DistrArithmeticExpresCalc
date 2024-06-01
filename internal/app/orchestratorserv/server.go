package orchestratorserv

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Stroevik3/DistrArithmeticExpresCalc/internal/app/orchestratorserv/model"
	"github.com/Stroevik3/DistrArithmeticExpresCalc/internal/app/orchestratorserv/store"
	"github.com/Stroevik3/DistrArithmeticExpresCalc/internal/app/orchestratorserv/store/memstore"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type server struct {
	mux    *http.ServeMux
	logger *logrus.Logger
	store  store.Store
}

func newServer(store store.Store) *server {
	s := &server{
		mux:    http.NewServeMux(),
		logger: logrus.New(),
		store:  store,
	}

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.logger.Debugf("started %s %s", r.Method, r.RequestURI)
	start := time.Now()
	s.mux.ServeHTTP(w, r)
	s.logger.Debugf(
		"completed in %v",
		time.Since(start),
	)
}

func (s *server) SetEpressionHandler(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugln("SetCalculateHandler")

		reqTxt, err := io.ReadAll(r.Body)
		if err != nil {
			s.logger.Error("ReadAll err - ", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var exp model.Expression
		err = json.Unmarshal(reqTxt, &exp)
		if err != nil {
			s.logger.Error("Unmarshal err - ", err.Error())
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		err = s.store.Expression().Add(&exp)
		if err != nil {
			s.logger.Error("Expression().Add err - ", err.Error())
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		go BreakExpressionIntoTasks(&exp, s.store)

		s.logger.Debugln("StatusCreated")
		w.WriteHeader(http.StatusCreated)
	})
}

func (s *server) GetEpressionsHandler(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugln("GetEpressionsHandler")

		var (
			expList []*model.Expression
			err     error
		)

		expList, err = s.store.Expression().GetList(true)
		if err != nil {
			s.logger.Error("Expression().GetList err - ", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		expListJsonByt, err := json.Marshal(&expList)
		if err != nil {
			s.logger.Error("Marshal err - ", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, "{\"expressions\":", string(expListJsonByt), "}")
	})
}

func (s *server) GetEpressionHandler(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugln("GetEpressionHandler")
		sId := r.PathValue("id")
		id, err := strconv.Atoi(sId)
		if err != nil {
			s.logger.Error("Atoi err - ", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		exp, err := s.store.Expression().Find(id)
		if err != nil {
			s.logger.Error("Find err - ", err.Error())
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		expResp := &model.Expression{
			Id:     exp.Id,
			Status: exp.Status,
			Result: exp.Result,
		}

		expJsonByt, err := json.Marshal(&expResp)
		if err != nil {
			s.logger.Error("Marshal err - ", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, "{\"expression\":", string(expJsonByt), "}")
	})
}

type TaskResp struct {
	Id         uuid.UUID `json:"id"`
	Exp        int       `json:"ExpId"`
	ArgOne     float64   `json:"arg1"`
	ArgTaskOne uuid.UUID `json:"argTaskIdOne"`
	ArgTwo     float64   `json:"arg2"`
	ArgTaskTwo uuid.UUID `json:"argTaskTwo"`
	Oper       string    `json:"operation"`
	OperTime   int       `json:"operation_time"`
	Prior      int       `json:"prior"`
	Status     string    `json:"status"`
	Result     float64   `json:"result"`
	TaskNext   uuid.UUID `json:"taskIdNext"`
}

func (s *server) GetTasksHandler(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugln("GetTasksHandler")

		var (
			taskList []*model.Task
			err      error
		)

		taskList, err = s.store.Task().GetList()
		if err != nil {
			s.logger.Error("Task().GetList err - ", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tasksResp := make([]*TaskResp, 0, len(taskList))
		getTaskId := func(ts *model.Task) uuid.UUID {
			var id uuid.UUID
			if ts != nil {
				return ts.Id
			}
			return id
		}

		for _, task := range taskList {
			tasksResp = append(tasksResp, &TaskResp{
				Id:         task.Id,
				Exp:        task.Exp.Id,
				ArgOne:     task.ArgOne,
				ArgTaskOne: getTaskId(task.ArgTaskOne),
				ArgTwo:     task.ArgTwo,
				ArgTaskTwo: getTaskId(task.ArgTaskTwo),
				Oper:       task.Oper,
				OperTime:   task.OperTime,
				Prior:      task.Prior,
				Status:     task.Status,
				Result:     task.Result,
				TaskNext:   getTaskId(task.TaskNext),
			})
		}

		taskListJsonByt, err := json.Marshal(&tasksResp)
		if err != nil {
			s.logger.Error("Marshal err - ", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, "{\"tasks\":", string(taskListJsonByt), "}")
	})
}

type TaskToCompleteResp struct {
	Id       uuid.UUID `json:"id"`
	Arg1     float64   `json:"arg1"`
	Arg2     float64   `json:"arg2"`
	Oper     string    `json:"operation"`
	OperTime int       `json:"operation_time"`
}

func (s *server) GetTaskToCompleteHandler(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugln("GetTaskToCompleteHandler")

		var (
			task *model.Task
			err  error
		)
		for {
			task, err = s.store.StackReadyTask().Pop()
			if err != nil {
				s.logger.Error("StackReadyTask().Pop() err - ", err.Error())
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			if task.Exp != nil && task.Exp.Status != memstore.EXP_STATUS_ERORR {
				break
			}
		}

		tasksResp := &TaskToCompleteResp{
			Id:       task.Id,
			Arg1:     task.ArgOne,
			Arg2:     task.ArgTwo,
			Oper:     task.Oper,
			OperTime: task.OperTime,
		}

		taskJsonByt, err := json.Marshal(&tasksResp)
		if err != nil {
			s.logger.Error("Marshal err - ", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		task.Status = memstore.TASK_STATUS_PROC
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, "{\"task\":", string(taskJsonByt), "}")
	})
}

type TaskResultReq struct {
	Id     uuid.UUID `json:"id"`
	Result float64   `json:"result"`
}

func (s *server) PostTaskResultHandler(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugln("PostTaskResultHandler")

		reqTxt, err := io.ReadAll(r.Body)
		if err != nil {
			s.logger.Error("ReadAll err - ", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var taskRes TaskResultReq
		err = json.Unmarshal(reqTxt, &taskRes)
		if err != nil {
			s.logger.Error("Unmarshal err - ", err.Error())
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		t, err := s.store.Task().Find(taskRes.Id)
		if err != nil {
			s.logger.Error("Unmarshal err - ", err.Error())
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if t.Status != memstore.TASK_STATUS_PROC {
			http.Error(w, "task in inappropriate status", http.StatusUnprocessableEntity)
			return
		}

		t.Result = taskRes.Result

		if t.TaskNext == nil {
			t.Exp.Result = strconv.FormatFloat(t.Result, 'f', 2, 64)
			t.Exp.Status = memstore.EXP_STATUS_RES_DETERM
		} else {
			nextTask := t.TaskNext
			if nextTask.ArgTaskOne != nil && nextTask.ArgTaskOne.Id == t.Id {
				nextTask.ArgOne = t.Result
				nextTask.ArgTaskOne = nil
			} else {
				nextTask.ArgTwo = t.Result
				nextTask.ArgTaskTwo = nil
			}
			if nextTask.ArgTaskOne == nil && nextTask.ArgTaskTwo == nil {
				err := s.store.StackReadyTask().Push(nextTask)
				if err != nil {
					nextTask.Exp.Status = memstore.EXP_STATUS_ERORR
					nextTask.Exp.Result = err.Error()
					return
				}
			}
		}
		t.Status = memstore.TASK_STATUS_RES_DETERM

		w.WriteHeader(http.StatusOK)
	})
}
