package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/MonoBear123/MyApi/back/model"
)

func main() {
	serverURL := "http://localhost:8080/setstatus" // Замените на правильный адрес вашего сервера
	id := rand.Uint64()
	out := model.Requert{Id: id, Status: "OK", Time: time.Now()}
	res, err := json.Marshal(out)
	if err != nil {
		return
	}
	for {
		// Формируем JSON с данными о статусе агента

		// Отправляем POST запрос на сервер
		response, err := http.Post(serverURL, "application/json", bytes.NewBuffer(res))
		if err != nil {
			fmt.Println("Error connecting to the server:", err)
		} else {

			fmt.Println("Server Status:", response.Status)
		}

		time.Sleep(5 * time.Second) // Пауза перед следующим запросом
	}
}
