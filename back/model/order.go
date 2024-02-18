package model

import (
	"time"

	shuntingYard "github.com/mgenware/go-shunting-yard"
)

type EXpression struct {
	Expression  string                   `json:"expression"`
	ExpressinID uint64                   `json:"expression_id"`
	ParsedEx    []*shuntingYard.RPNToken `json:"parsedex"`
	Result      []*shuntingYard.RPNToken `json:"result"`
}

type Requert struct {
	Id            uint64    `json:"id"`
	Status        string    `json:"status"`
	Time          time.Time `json:"time"`
	NumOfWorkers  int       `json:"numofworkers"`
	MaxNumWorkers int       `json:"maxnumworkers"`
	
}
type Config struct {
	Plus           int `json:"+"`
	Minus          int `json:"-"`
	Division       int `json:"/"`
	Multiplication int `json:"*"`
	Construction   int `json:"^"`
	MaxGorutines   int `json:"gorutines"`
}
