package model

type Expression struct {
	Id     int    `json:"id"`
	Val    string `json:"expression,omitempty"`
	Status string `json:"status"`
	Result string `json:"result"`
}
