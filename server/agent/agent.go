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
	"github.com/redis/go-redis/v9"
)

const (
	heartbeatPeriod = 15 * time.Second
)

var (
	mutex         sync.Mutex
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

	if err := clientRedis.Ping(context.Background()).Err(); err != nil {
		fmt.Print("не подлючен к редис")
	}
	fmt.Println("Успешное подключение к Redis")

	for {
		expression, err := LockQueue("my_queue", clientRedis)
		if err != nil {
			fmt.Println(err)
			continue
		}
		ConfigJson, err := clientRedis.Get(context.Background(), "config").Result()
		if err != nil {
			fmt.Print("не получил конфига ")
		}
		config2 := model.Config{}

		err = json.Unmarshal([]byte(ConfigJson), &config2)
		if err != nil {
			fmt.Print("не подлючен к редис")
		}
		// Получаем значение из очереди
		fmt.Println("значение получено")
		//mutex.Lock()

		//mutex.Unlock()
		fmt.Println(expression.Num1, " ", expression.Num2)
		workerSemaphore <- struct{}{}

		fmt.Println("после канала ", expression.Id)
		go func(expression SubEx, clientRedis *redis.Client) {

			activeWorkers++

			var res float64
			var Error string
			switch expression.Operator {
			case "/":
				fmt.Println("запустился на делении")
				if expression.Num2 == 0.0 {
					Error = "err"
				}
				time.Sleep(time.Second * time.Duration(config2.Division))
				res = expression.Num1 / expression.Num2

			case "*":
				fmt.Println("запустился на произведении")
				time.Sleep(time.Second * time.Duration(config2.Multiplication))
				res = expression.Num1 * expression.Num2

			case "-":
				fmt.Println("запустился на минусе")
				time.Sleep(time.Second * time.Duration(config2.Minus))

				res = expression.Num1 - expression.Num2

			case "+":
				fmt.Println("запустился на плюсе")
				fmt.Println(expression.Num1, " ", expression.Num2)
				time.Sleep(time.Second * time.Duration(config2.Plus))
				res = expression.Num1 + expression.Num2
			case "^":
				time.Sleep(time.Second * time.Duration(config2.Construction))

				res = math.Pow(expression.Num1, expression.Num2)
			default:
				fmt.Println("не найден оператор")
			}

			out, err := json.Marshal(Result{
				Res:   res,
				Index: expression.Index,
				Error: Error,
			})
			if err != nil {
				fmt.Print("не удалось замарщалить результат")
			}
			fmt.Println(res)

			err = clientRedis.LPush(context.Background(), expression.Id, string(out)).Err()
			if err != nil {
				fmt.Printf("Ошибка при отправке данных в очередь: %v\n", err)
			} else {
				fmt.Println("Данные успешно отправлены в очередь")
			}

			mutex.Lock()
			activeWorkers--
			mutex.Unlock()

			<-workerSemaphore

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

	result, err := client.BLPop(ctx, 0, queueName).Result()
	if err != nil {
		return SubEx{}, fmt.Errorf("ошибка при ожидании значения из очереди: %v", err)
	}
	expression := SubEx{}

	err = json.Unmarshal([]byte(result[1]), &expression)
	if err != nil {
		return SubEx{}, fmt.Errorf("ошибка распаковки JSON: %v", err)
	}

	return expression, nil
}
