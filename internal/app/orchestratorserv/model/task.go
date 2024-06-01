package model

import "github.com/google/uuid"

type Task struct {
	Id         uuid.UUID   `json:"id"`
	Exp        *Expression `json:"-"`
	ArgOne     float64     `json:"arg1"`
	ArgTaskOne *Task       `json:"-"`
	ArgTwo     float64     `json:"arg2"`
	ArgTaskTwo *Task       `json:"-"`
	Oper       string      `json:"operation"`
	OperTime   int         `json:"operation_time"`
	Prior      int         `json:"-"`
	Status     string      `json:"-"`
	Result     float64     `json:"result"`
	TaskNext   *Task       `json:"-"`
}
