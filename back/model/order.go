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
	Plus           int `json:"+" env-default:"100"`
	Minus          int `json:"-" env-default:"100"`
	Division       int `json:"/" env-default:"100"`
	Multiplication int `json:"*" env-default:"100"`
	Construction   int `json:"^" env-default:"100"`
	MaxGorutines   int `json:"gorutines" env-default:"10"`
}
