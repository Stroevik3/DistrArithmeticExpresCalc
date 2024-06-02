package model

import (
	"time"

	"github.com/google/uuid"
)

type TaskGetResp struct {
	Id       uuid.UUID     `json:"id"`
	ArgOne   float64       `json:"arg1"`
	ArgTwo   float64       `json:"arg2"`
	Oper     string        `json:"operation"`
	OperTime time.Duration `json:"operation_time"`
}
