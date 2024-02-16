package expression

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"fmt"

	"strings"

	"github.com/MonoBear123/MyApi/back/model"
	shuntingYard "github.com/mgenware/go-shunting-yard"
)

func ParseExpression(expr string) ([]*shuntingYard.RPNToken, error) {
	infixTokens, err := shuntingYard.Scan(expr)
	if err != nil {
		return nil, err
	}

	postfixTokens, err := shuntingYard.Parse(infixTokens)
	if err != nil {
		return nil, err
	}

	return postfixTokens, nil
}
func (r *RedisRepo) Distribution(expression []*shuntingYard.RPNToken, id string) ([]*shuntingYard.RPNToken, error) {
	per, err := r.Client.Get(context.Background(), "config").Result()
	if err != nil {
		return nil, fmt.Errorf("ошибка при получения конфига для парсинга")
	}
	config := model.Config{}
	err = json.Unmarshal([]byte(per), &config)
	if err != nil {
		return nil, fmt.Errorf("ошибка вывода из джейсона")
	}
	maxTime := max(config.Construction, config.Minus, config.Division, config.Plus, config.Multiplication)
	for len(expression) != 1 {
		for index := 0; index < len(expression)-3; index++ {
			if !strings.ContainsAny(fmt.Sprintf("%v", expression[index].Value), "+-/*^") && !strings.ContainsAny(fmt.Sprintf("%v", expression[index+1].Value), "+-/*^") && strings.ContainsAny(fmt.Sprintf("%v", expression[index+2].Value), "+-/*^") {

				err := r.EnqueueMessage("my_queue", []interface{}{expression[index].Value, expression[index+1].Value, expression[index+2].Value, id, index}, maxTime)
				if err != nil {
					fmt.Print(err)
				}
				log.Printf("отправлено в очередь %d %d %d", expression[index].Value, expression[index+1].Value, expression[index+2].Value)
				expression[index].Value = "."
				expression[index+1].Value = "."
				expression[index+2].Value = "."

			}

		}
		//тут будет тайслип на время выполнения операций
		time.Sleep(time.Duration(maxTime+2) * time.Second)
		for {
			newEX, err := r.DequeueMessage(id)
			if err != nil {
				break
			}
			if newEX[2] == "err" {
				return nil, fmt.Errorf("встречено деление на ноль")
			}

			expression[newEX[1].(int)].Value = newEX[0]
			expression = append(expression[:newEX[1].(int)], expression[newEX[1].(int)+2:]...)

			r.Client.JSONSet(context.Background(), id, ".", fmt.Sprintf(`{"parsedex":"%v"}`, expression))
		}

	}

	return expression, nil
}
func (r *RedisRepo) EnqueueMessage(name string, subEx []interface{}, maxTime int) error {
	ctx := context.Background()
	err := r.Client.LPush(ctx, name, subEx)
	if err != nil {
		return fmt.Errorf("ошибка в очереди")

	}
	r.Client.Expire(context.Background(), name, time.Duration(maxTime)*time.Minute)
	return nil
}
func (r *RedisRepo) DequeueMessage(name string) ([]interface{}, error) {
	result, err := r.Client.LPop(context.Background(), name).Result()
	if err != nil {
		return nil, fmt.Errorf("ошибка в очереди: %v", err)
	}

	var expression []interface{}
	err = json.Unmarshal([]byte(result), &expression)
	if err != nil {
		return nil, fmt.Errorf("ошибка распаковки JSON: %v", err)
	}

	return expression, nil
}
