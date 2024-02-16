package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/MonoBear123/MyApi/back/model"
	"github.com/redis/go-redis/v9"
)

const (
	heartbeatPeriod = 15 * time.Second
)

var (
	mutex         sync.Mutex
	activeWorkers int
)

func main() {
	client := &http.Client{}
	options := &redis.Options{
		Addr: "redis:6379", // замените на реальный адрес

	}
	clientRedis := redis.NewClient(options)
	id := rand.Uint64()
	ConfigJson, err := clientRedis.Get(context.Background(), "config").Result()
	if err != nil {
		fmt.Print("не получил конфига ")
	}
	config := model.Config{}
	err = json.Unmarshal([]byte(ConfigJson), &config)
	if err != nil {
		fmt.Print("не подлючен к редис")
	}
	workerSemaphore := make(chan struct{}, config.MaxGorutines)
	go sendHeartbeat(client, "http://server:8041/setstatus", id, config.MaxGorutines)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := clientRedis.Ping(ctx).Err(); err != nil {
		fmt.Print("не подлючен к редис")
	}

	fmt.Println("Успешное подключение к Redis")
	for {
		expression, err := LockQueue("my_queue", clientRedis)
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
			var Error string
			switch expression[2] {
			case "/":
				if expression[1] == 0 {
					Error = "err"
				}
				time.Sleep(time.Second * time.Duration(config.Division))
				num1 := expression[1].(float64)
				num2 := expression[0].(float64)
				res = num1 / num2
			case "*":
				time.Sleep(time.Second * time.Duration(config.Multiplication))
				num1 := expression[1].(float64)
				num2 := expression[0].(float64)
				res = num1 * num2

			case "-":
				time.Sleep(time.Second * time.Duration(config.Minus))
				num1 := expression[1].(float64)
				num2 := expression[0].(float64)
				res = num1 - num2

			case "+":
				time.Sleep(time.Second * time.Duration(config.Plus))
				num1 := expression[1].(float64)
				num2 := expression[0].(float64)
				res = num1 + num2
			case "^":
				time.Sleep(time.Second * time.Duration(config.Construction))
				num1 := expression[1].(float64)
				num2 := expression[0].(float64)
				res = math.Pow(num1, num2)
			}

			log.Print(res)
			clientRedis.LPush(context.Background(), expression[4].(string), []interface{}{res, expression[4], Error})

		}(expression)

	}

}
func sendHeartbeat(client *http.Client, url string, id uint64, MaxWorkers int) {
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

func LockQueue(queueName string, client *redis.Client) ([]interface{}, error) {
	ctx := context.Background()
	lockKey := queueName + "_lock"
	lockValue := "locked"
	lockSet, err := client.SetNX(ctx, lockKey, lockValue, time.Second).Result()
	if err != nil {
		return nil, fmt.Errorf("ошибка при установке блокировки: %v", err)
	}

	if !lockSet {
		// Блокировка уже установлена, пропускаем задачу
		return nil, fmt.Errorf("12")
	}

	defer func(clientRedis *redis.Client) {
		// Снимаем блокировку после выполнения задачи
		clientRedis.Del(ctx, lockKey)
	}(client)
	result, err := client.LPop(ctx, queueName).Result()
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
