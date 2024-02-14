package model

import "time"

type EXpression struct {
	Expression  string  `json:"expression"`
	ExpressinID uint64  `json:"expression_id"`
	Tree        *Node   `json:"tree"`
	Num         float64 `json:"num"`
}

type Node struct {
	Left     *Node
	Right    *Node
	Operator string
	Value    float64
}

type Requert struct {
	Id     uint64    `json:"id"`
	Status string    `json:"status"`
	Time   time.Time `json:"time"`
}
