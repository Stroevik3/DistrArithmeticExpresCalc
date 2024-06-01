package memstore

import (
	"errors"
	"regexp"
	"strings"
	"sync"

	"github.com/Stroevik3/DistrArithmeticExpresCalc/internal/app/orchestratorserv/model"
	"github.com/Stroevik3/DistrArithmeticExpresCalc/internal/app/orchestratorserv/store"
)

const (
	EXP_STATUS_CREATE     string = "CREATE"         // Создано
	EXP_STATUS_PROC       string = "PROCESSING"     // Обрабатывается
	EXP_STATUS_ERORR      string = "ERORR"          // Ошибка
	EXP_STATUS_RES_DETERM string = "RES_DETERMINED" // Результат определен
)

type ExpressionRepository struct {
	mu          sync.Mutex
	store       *Store
	expressions map[int]*model.Expression
}

func ValidateExp(expVal string) error {
	val := strings.ReplaceAll(expVal, "+", "")
	val = strings.ReplaceAll(val, "-", "")
	val = strings.ReplaceAll(val, "/", "")
	val = strings.ReplaceAll(val, "*", "")
	val = strings.ReplaceAll(val, "(", "")
	val = strings.ReplaceAll(val, ")", "")
	val = strings.ReplaceAll(val, ".", "")
	var re = regexp.MustCompile(`^[0-9]+$`)
	if !re.MatchString(val) {
		return errors.New("invalid expression")
	}
	return nil
}

func (r *ExpressionRepository) Add(e *model.Expression) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	e.Status = EXP_STATUS_CREATE
	e.Result = ""
	_, ok := r.expressions[e.Id]
	if ok {
		return store.ErrRecordDubleId
	}
	if err := ValidateExp(e.Val); err != nil {
		return err
	}
	r.expressions[e.Id] = e
	return nil
}

func (r *ExpressionRepository) GetList(clearVl bool) ([]*model.Expression, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	newExpressions := make([]*model.Expression, 0, len(r.expressions))
	for _, val := range r.expressions {
		newExpressions = append(newExpressions, &model.Expression{
			Id: val.Id,
			Val: func(cl bool, vl string) string {
				if clearVl {
					return ""
				} else {
					return vl
				}
			}(clearVl, val.Val),
			Status: val.Status,
			Result: val.Result,
		})
	}
	return newExpressions, nil
}

func (r *ExpressionRepository) Find(id int) (*model.Expression, error) {
	expression, ok := r.expressions[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return expression, nil
}
