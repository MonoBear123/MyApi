package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"

	//"sync"
	"time"

	"github.com/MonoBear123/MyApi/back/model"
	"github.com/redis/go-redis/v9"
)

const (
	heartbeatPeriod = 15 * time.Second
)

var (
	//mutex         sync.Mutex
	activeWorkers int
)

type Result struct {
	Res   float64 `json:"res"`
	Index int     `json:"index"`
	Error string  `json:"error"`
}
type SubEx struct {
	Num1     float64 `json:"num1"`
	Num2     float64 `json:"num2"`
	Operator string  `json:"operator"`
	Id       string  `json:"id"`
	Index    int     `json:"index"`
}

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
		fmt.Println("значение получено")
		//mutex.Lock()

		activeWorkers++
		fmt.Println("доп канала ", expression.Num1, expression.Num2)
		//mutex.Unlock()
		fmt.Println("мьютексы пройдены?")
		workerSemaphore <- struct{}{}
		fmt.Println("после канала ", expression.Num1, expression.Num2)
		go func(expression SubEx, clientRedis *redis.Client) {
			defer func() {

				<-workerSemaphore
				//mutex.Lock()
				activeWorkers--
				//mutex.Unlock()
			}()
			var res float64
			var Error string
			switch expression.Operator {
			case "/":
				fmt.Print("запустился на делении")
				if expression.Num2 == 0 {
					Error = "err"
				}
				time.Sleep(time.Second * time.Duration(config.Division))
				res = expression.Num1 / expression.Num2

			case "*":
				fmt.Print("запустился на произведении")
				time.Sleep(time.Second * time.Duration(config.Multiplication))
				res = expression.Num1 * expression.Num2

			case "-":
				fmt.Print("запустился на минусе")
				time.Sleep(time.Second * time.Duration(config.Minus))

				res = expression.Num1 - expression.Num2

			case "+":
				fmt.Print("запустился на плюсе")
				time.Sleep(time.Second * time.Duration(config.Plus))
				res = expression.Num1 + expression.Num2
			case "^":
				time.Sleep(time.Second * time.Duration(config.Construction))

				res = math.Pow(expression.Num1, expression.Num2)
			default:
				fmt.Print("не найден оператор")
			}

			out, err := json.Marshal(Result{
				Res:   res,
				Index: expression.Index,
				Error: Error,
			})
			if err != nil {
				fmt.Print("не удалось замарщалить результат")
			}
			fmt.Print(res)
			clientRedis.LPush(context.Background(), expression.Id, out)

		}(expression, clientRedis)

	}

}
func sendHeartbeat(client *http.Client, url string, id uint64, MaxWorkers int) {

	for {
		//mutex.Lock()

		//mutex.Unlock()
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

func LockQueue(queueName string, client *redis.Client) (SubEx, error) {
	ctx := context.Background()
	lockKey := queueName + "_lock"
	lockValue := "locked"
	lockSet, err := client.SetNX(ctx, lockKey, lockValue, time.Second).Result()
	if err != nil {
		return SubEx{}, fmt.Errorf("ошибка при установке блокировки: %v", err)
	}

	if !lockSet {
		// Блокировка уже установлена, пропускаем задачу
		return SubEx{}, fmt.Errorf("12")
	}

	defer func(clientRedis *redis.Client) {
		// Снимаем блокировку после выполнения задачи
		clientRedis.Del(ctx, lockKey)
	}(client)

	result, err := client.LPop(ctx, queueName).Result()
	if err != nil {
		return SubEx{}, fmt.Errorf("ошибка в очереди: %v", err)
	}

	expression := SubEx{}

	err = json.Unmarshal([]byte(result), &expression)
	if err != nil {
		return SubEx{}, fmt.Errorf("ошибка распаковки JSON: %v", err)
	}

	return expression, nil
}
