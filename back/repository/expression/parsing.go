package expression

import (
	"context"
	"encoding/json"

	"fmt"

	"strings"

	"github.com/MonoBear123/MyApi/back/model"
	shuntingYard "github.com/mgenware/go-shunting-yard"
)

type SubEx struct {
	Num1     float64 `json:"num1"`
	Num2     float64 `json:"num2"`
	Operator string  `json:"operator"`
	Id       string  `json:"id"`
	Index    int     `json:"index"`
}
type Result struct {
	Res   float64 `json:"res"`
	Index int     `json:"index"`
	Error string  `json:"error"`
}

func ParseExpression(expr string) ([]*shuntingYard.RPNToken, error) {
	infixTokens, err := shuntingYard.Scan(expr)
	if err != nil {
		return nil, err
	}

	postfixTokens, err := shuntingYard.Parse(infixTokens)
	if err != nil {
		return nil, err
	}
	fmt.Println("прошло обработку")
	return postfixTokens, nil
}
func (r *RedisRepo) Distribution(expression []*shuntingYard.RPNToken, id uint64, ex string) ([]*shuntingYard.RPNToken, error) {
	per, err := r.Client.Get(context.Background(), "config").Result()
	if err != nil {
		return nil, fmt.Errorf("ошибка при получения конфига для парсинга")
	}
	config := model.Config{}
	err = json.Unmarshal([]byte(per), &config)
	if err != nil {
		return nil, fmt.Errorf("ошибка вывода из джейсона")
	}
	fmt.Println(fmt.Sprint(id))
	colex := 0
	maxTime := max(config.Construction, config.Minus, config.Division, config.Plus, config.Multiplication)
	fmt.Println("начало обработки выражения  ")
	for len(expression) != 1 {
		for index := 0; index < len(expression)-2; index++ {

			if !strings.ContainsAny(fmt.Sprint(expression[index].Value), "+-/*^.") && !strings.ContainsAny(fmt.Sprint(expression[index+1].Value), "+-/*^.") && strings.ContainsAny(fmt.Sprint(expression[index+2].Value), "+-/*^") {
				err := r.EnqueueMessage("my_queue", SubEx{
					Num1:     float64(expression[index].Value.(int)),
					Num2:     float64(expression[index+1].Value.(int)),
					Operator: expression[index+2].Value.(string),
					Id:       "qeue:" + fmt.Sprint(id),
					Index:    index,
				}, maxTime)
				if err != nil {
					fmt.Println("НЕ ОТПРАВИЛОСЬ В ОЧЕРЕДЬ ", err)
				}
				colex++
				fmt.Printf("отправлено в очередь %d %d %d", expression[index].Value, expression[index+1].Value, expression[index+2].Value)
				expression[index].Value = "."
				expression[index+1].Value = "."
				expression[index+2].Value = "."

			}

		}
		//тут будет тайслип на время выполнения операций

		for colex != 0 {
			fmt.Println(colex)
			fmt.Println("ждет отправления")
			newEX, err := r.DequeueMessage(fmt.Sprint(id))
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println(colex)
			colex--
			fmt.Println(colex)
			if newEX.Error == "err" {
				return nil, fmt.Errorf("встречено деление на ноль")
			}
			expression[newEX.Index].Value = newEX.Res
			if len(expression) < 4 {
				expression = expression[:newEX.Index+1]
			} else {
				expression = append(expression[:newEX.Index+1], expression[newEX.Index+3:]...)
			}

			out := model.EXpression{
				Expression:  ex,
				ExpressinID: id,
				ParsedEx:    expression,
			}
			r.Client.Set(context.Background(), fmt.Sprint(id), out, 0)
		}

	}

	return expression, nil
}
func (r *RedisRepo) EnqueueMessage(name string, subEx SubEx, maxTime int) error {
	ctx := context.Background()
	res, err := json.Marshal(subEx)
	if err != nil {
		return err
	}
	err = r.Client.LPush(ctx, name, string(res)).Err()
	if err != nil {
		fmt.Println("не попал в очередб")

		return fmt.Errorf("ошибка в очереди")

	}

	return nil
}
func (r *RedisRepo) DequeueMessage(name string) (Result, error) {
	result, err := r.Client.BLPop(context.Background(), 0, "qeue:"+name).Result()
	if err != nil {
		return Result{}, fmt.Errorf("ошибка при ожидании значения из очереди: %v", err)
	}
	fmt.Println("принял значение в очереди аосле ")

	var expression Result
	err = json.Unmarshal([]byte(result[1]), &expression)
	if err != nil {
		return Result{}, fmt.Errorf("ошибка распаковки JSON: %v", err)
	}

	return expression, nil
}
