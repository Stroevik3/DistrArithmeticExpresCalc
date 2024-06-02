package orchestratorserv

import (
	"strconv"
	"time"

	"github.com/Stroevik3/DistrArithmeticExpresCalc/internal/app/orchestratorserv/model"
	"github.com/Stroevik3/DistrArithmeticExpresCalc/internal/app/orchestratorserv/store"
	"github.com/Stroevik3/DistrArithmeticExpresCalc/internal/app/orchestratorserv/store/memstore"
)

const (
	OPER_SYMB_ADDIT string = "+"
	OPER_SYMB_SUBTR string = "-"
	OPER_SYMB_MULTP string = "*"
	OPER_SYMB_DIVIS string = "/"

	BRACKET_LEFT  string = "("
	BRACKET_RIGHT string = ")"
	DOT_SYMBOK    string = "."

	TIME_ADDITION_MS        time.Duration = (7000 * time.Millisecond)  // время выполнения операции сложения в милисекундах
	TIME_SUBTRACTION_MS     time.Duration = (1700 * time.Millisecond)  // время выполнения операции вычитания в милисекундах
	TIME_MULTIPLICATIONS_MS time.Duration = (26000 * time.Millisecond) // время выполнения операции умножения в милисекундах
	TIME_DIVISIONS_MS       time.Duration = (38000 * time.Millisecond) // время выполнения операции деления в милисекундах

	// Чем больше приоритет тем он выше
	PRIORITY_ONE int = 1
	PRIORITY_TWO int = 2
)

func GetTimeMsOp(op string) time.Duration {
	var res time.Duration
	switch op {
	case OPER_SYMB_ADDIT:
		res = TIME_ADDITION_MS
	case OPER_SYMB_SUBTR:
		res = TIME_SUBTRACTION_MS
	case OPER_SYMB_MULTP:
		res = TIME_MULTIPLICATIONS_MS
	case OPER_SYMB_DIVIS:
		res = TIME_DIVISIONS_MS
	}
	return res
}

func GetPriorityOp(op string) int {
	var res int
	switch {
	case op == OPER_SYMB_ADDIT || op == OPER_SYMB_SUBTR:
		res = PRIORITY_ONE
	case op == OPER_SYMB_MULTP || op == OPER_SYMB_DIVIS:
		res = PRIORITY_TWO
	}
	return res
}

func BreakExpressionIntoTasks(exp *model.Expression, st store.Store) {
	sVl := exp.Val
	var (
		op, sArg   string
		arg        float64
		deltaPrior int
		tsPrev     *model.Task
		ts         *model.Task
		tsNext     *model.Task
	)
	for i := 0; i < len(sVl); i++ {
		sSymb := string(sVl[i])

		_, err := strconv.Atoi(sSymb)
		if sSymb != DOT_SYMBOK && err != nil {
			if sSymb == BRACKET_LEFT {
				deltaPrior += 2
				continue
			}
			if sSymb == BRACKET_RIGHT {
				deltaPrior -= 2
				if deltaPrior < 0 {
					deltaPrior = 0
				}
				continue
			}
			if sArg == "" {
				sArg = sSymb
				continue
			}
			op = sSymb
			arg, err = strconv.ParseFloat(sArg, 64)
			if err != nil {
				sArg = op
				op = ""
				continue
			}
			ts = &model.Task{
				Exp:      exp,
				Oper:     op,
				OperTime: GetTimeMsOp(op),
				Prior:    (GetPriorityOp(op) + deltaPrior),
			}
			if tsPrev == nil {
				ts.ArgOne = arg
				tsPrev = ts
			} else {
				if tsPrev.Prior >= ts.Prior {
					tsPrev.ArgTwo = arg
					if tsPrev.TaskNext == nil {
						tsPrev.TaskNext = ts
						ts.ArgTaskOne = tsPrev
					} else {
						if tsPrev.TaskNext.Prior < ts.Prior {
							tsNext = tsPrev.TaskNext
							tsPrev.TaskNext = ts
							ts.ArgTaskOne = tsPrev
							ts.TaskNext = tsNext
							tsNext.ArgTaskTwo = ts
						} else {
							if tsPrev.TaskNext.TaskNext == nil {
								tsPrev.TaskNext.TaskNext = ts
								ts.ArgTaskOne = tsNext
							} else {
								tsPrev.TaskNext.TaskNext.TaskNext = ts
								ts.ArgTaskOne = tsPrev.TaskNext.TaskNext
							}
						}
					}
					err := st.Task().Add(tsPrev)
					if err != nil {
						exp.Status = memstore.EXP_STATUS_ERORR
						exp.Result = err.Error()
						return
					}
					if tsPrev.ArgTaskOne == nil && tsPrev.ArgTaskTwo == nil {
						err := st.StackReadyTask().Push(tsPrev)
						if err != nil {
							exp.Status = memstore.EXP_STATUS_ERORR
							exp.Result = err.Error()
							return
						}
						tsPrev.Status = memstore.TASK_STATUS_READY
					}
				} else {
					tsPrev.ArgTaskTwo = ts
					ts.ArgOne = arg
					ts.TaskNext = tsPrev
					err := st.Task().Add(tsPrev)
					if err != nil {
						exp.Status = memstore.EXP_STATUS_ERORR
						exp.Result = err.Error()
						return
					}
				}
				tsPrev = ts
			}
			sArg = ""

		} else {
			sArg += sSymb
		}
	}
	if sArg != "" {
		arg, err := strconv.ParseFloat(sArg, 64)
		if err == nil {
			ts.ArgTwo = arg
			err := st.Task().Add(ts)
			if err != nil {
				exp.Status = memstore.EXP_STATUS_ERORR
				exp.Result = err.Error()
				return
			}
			if ts.ArgTaskOne == nil && ts.ArgTaskTwo == nil {
				err := st.StackReadyTask().Push(ts)
				if err != nil {
					exp.Status = memstore.EXP_STATUS_ERORR
					exp.Result = err.Error()
					return
				}
				ts.Status = memstore.TASK_STATUS_READY
			}

		}
	}
	exp.Status = memstore.EXP_STATUS_PROC
}
