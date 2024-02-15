package expression

import (
	"context"
	"encoding/json"
	"time"

	"fmt"
	"log"
	"strings"

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
func (r *RedisRepo) Distribution(expression []*shuntingYard.RPNToken, id string) {
	for len(expression) != 1 {
		for index := 0; index < len(expression)-3; index++ {
			if !strings.ContainsAny(fmt.Sprintf("%v", expression[index].Value), "+-/*^") && !strings.ContainsAny(fmt.Sprintf("%v", expression[index+1].Value), "+-/*^") && strings.ContainsAny(fmt.Sprintf("%v", expression[index+2].Value), "+-/*^") {

				err := r.EnqueueMessage("my_queue", []interface{}{expression[index].Value, expression[index+1].Value, expression[index+2].Value, id, index})
				if err != nil {
					log.Fatal(err)
				}
				expression[index].Value = "."
				expression[index+1].Value = "."
				expression[index+2].Value = "."

			}

		}
		//тут будет тайслип на время выполнения операций
		time.Sleep(100 * time.Second)
		for {
			newEX, err := r.DequeueMessage(id)
			if err != nil {
				break
			}
			expression[newEX[1].(int)].Value = newEX[0]
			expression = append(expression[:newEX[1].(int)], expression[newEX[1].(int)+2:]...)

			r.Client.JSONSet(context.Background(), id, ".", fmt.Sprintf(`{"parsedex":"%v"}`, expression))
		}

	}

	r.Client.JSONSet(context.Background(), id, ".", fmt.Sprintf(`{"result":"%v"}`, expression))
}
func (r *RedisRepo) EnqueueMessage(name string, subEx []interface{}) error {
	ctx := context.Background()
	err := r.Client.LPush(ctx, name, subEx)
	if err != nil {
		return fmt.Errorf("ошибка в очереди")
	}
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
