package model

import (
	"time"
)

type EXpression struct {
	Expression  string     `json:"expression"`
	ExpressinID uint64     `json:"expression_id"`
	CreatedEX   *time.Time `json:"created_ex"`
	WaitingEX   *time.Time `json:"waiting_ex"`
	CompelEx    *time.Time `json:"compel_ex"`
}
