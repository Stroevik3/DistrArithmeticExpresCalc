package model

import "github.com/google/uuid"

type TaskReq struct {
	Id     uuid.UUID `json:"id"`
	Result float64   `json:"result"`
}
