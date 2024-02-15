package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/MonoBear123/MyApi/back/model"
	"github.com/MonoBear123/MyApi/back/repository/expression"
)

type Req struct {
	Repo *expression.RedisRepo
}

const (
	MaxWorkers      = 8
	heartbeatPeriod = 15 * time.Second
)

var (
	workerSemaphore = make(chan struct{}, MaxWorkers)
	mutex           sync.Mutex
	activeWorkers   int
)

func (r *Req) main() {
	client := &http.Client{}
	id := rand.Uint64()
	go sendHeartbeat(client, "http://server:8041/setstatus", id)

	for {
		expression,err := r.LockQueue("my_queue")
		if err != nil {
			continue
		}
		// Получаем значение из очереди
		
		mutex.Lock()
		activeWorkers++
		mutex.Unlock()
		workerSemaphore <- struct{}{}

		go func(expression []interface{}) {
			defer func() {

				<-workerSemaphore
				mutex.Lock()
				activeWorkers--
				mutex.Unlock()
			}()
			var res float64
			switch expression[2] {
			case "/":
				if expression[1] == 0 {
					return
				}
				num1 := expression[1].(float64)
				num2 := expression[0].(float64)
				res = num1 / num2
			case "*":
				num1 := expression[1].(float64)
				num2 := expression[0].(float64)
				res = num1 * num2

			case "-":
				num1 := expression[1].(float64)
				num2 := expression[0].(float64)
				res = num1 - num2

			case "+":
				num1 := expression[1].(float64)
				num2 := expression[0].(float64)
				res = num1 / num2
			case "^":
				num1 := expression[1].(float64)
				num2 := expression[0].(float64)
				res = math.Pow(num1, num2)
			}
			time.Sleep(100 * time.Second)
			r.Repo.Client.LPush(context.Background(), expression[4].(string), []interface{}{res,expression[4]})
		}(expression)

	}

}
func sendHeartbeat(client *http.Client, url string, id uint64) {
	for {
		res, err := json.Marshal(model.Requert{Id: id, Status: "OK", Time: time.Now(),
			NumOfWorkers: activeWorkers, MaxNumWorkers: MaxWorkers})
		if err != nil {
			return
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(res))
		req.Header.Set("Content-Type", "application/json")

		if err != nil {
			return
		}

		resp, err := client.Do(req)

		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		time.Sleep(heartbeatPeriod)

	}
}

func (r *Req) LockQueue(queueName string) ([]interface{},error) {
	ctx := context.Background()
	lockKey := queueName + "_lock"
	lockValue := "locked"
	lockSet, err := r.Repo.Client.SetNX(ctx, lockKey, lockValue, time.Second).Result()
	if err != nil {
		return nil,fmt.Errorf("ошибка при установке блокировки: %v", err)
	}

	if !lockSet {
		// Блокировка уже установлена, пропускаем задачу
		return nil,fmt.Errorf("12")
	}

	defer func() {
		// Снимаем блокировку после выполнения задачи
		r.Repo.Client.Del(ctx, lockKey)
	}()
	result, err := r.Repo.Client.LPop(ctx, queueName).Result()
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
