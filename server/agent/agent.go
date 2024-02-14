package main

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/MonoBear123/MyApi/back/model"
)

func main() {
	client := &http.Client{}
	id := rand.Uint64()

	for i := 0; i < 5; i++ {
		res, err := json.Marshal(model.Requert{Id: id, Status: "OK", Time: time.Now()})
		req, err := http.NewRequest("POST", "http://server:8041/setstatus", bytes.NewBuffer(res))
		req.Header.Set("Content-Type", "application/json")

		if err != nil {
			return
		}

		resp, err := client.Do(req)

		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		time.Sleep(3 * time.Second)
	}
}
