package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"

	"math/rand"
	"net/http"
	"time"

	"github.com/MonoBear123/MyApi/back/model"
	"github.com/MonoBear123/MyApi/back/repository/expression"
)

type Expression struct {
	Repo *expression.RedisRepo
}

func (o *Expression) SetExpression(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Expression1 string `json:"expression"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	
	parsedExpression, err := expression.ParseExpression(body.Expression1)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Выражение не прошло валидацию"))
		return
	}
	id := rand.Uint64()
	out := model.EXpression{
		Expression:  body.Expression1,
		ExpressinID: id,
		ParsedEx:    parsedExpression,
	}

	err = o.Repo.Insert(r.Context(), out)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		w.Write([]byte("ошибка в базе данных"))
		return
	}

	outres, err := o.Repo.Distribution(parsedExpression, id, body.Expression1)
	if err != nil {
		o.Repo.DeleteAgent(fmt.Sprint(id))
		w.Write([]byte("встречено деление на 0"))
		return
	}



	w.Write([]byte(out.Expression + "=" + fmt.Sprint(outres[0].Value)))
	o.Repo.DeleteExpression(fmt.Sprint(id))

}

func (o *Expression) SetAgentStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var Newagent model.Requert
		err := json.NewDecoder(r.Body).Decode(&Newagent)

		if err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}

		o.Repo.AgentInsert(r.Context(), Newagent)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Received heartbeat"))
		return
	}
}
func (o *Expression) GetAgentStatus(w http.ResponseWriter, r *http.Request) {
	ctx := make(map[string]interface{})
	status := o.Repo.AgentALLFind(r.Context())
	for i, onestatus := range status {
		timeOld := time.Since(onestatus.Time)

		if timeOld > time.Second*30 {
			ctx[fmt.Sprintf("Агент %d", i+1)] = model.Requert{
				Status:        "Отключился",
				MaxNumWorkers: onestatus.MaxNumWorkers,
				NumOfWorkers:  onestatus.NumOfWorkers,
			}
			o.Repo.DeleteAgent(fmt.Sprint(onestatus.Id))
		} else {
			ctx[fmt.Sprintf("Агент %d", i+1)] = model.Requert{
				Status:        "ОК",
				MaxNumWorkers: onestatus.MaxNumWorkers,
				NumOfWorkers:  onestatus.NumOfWorkers,
			}
		}
	}
	t, _ := template.ParseFiles("front/pages/redis.html")
	t.Execute(w, ctx)
}
func (o *Expression) UpdateConfigHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	currentConfig := model.Config{}
	// Парсим значения из формы
	if err := json.NewDecoder(r.Body).Decode(&currentConfig); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}
	// Обновляем значения в конфигурации

	jsonConfig, err := json.Marshal(currentConfig)
	if err != nil {
		w.Write([]byte("ошибка в чтении конфига"))
	}

	err = o.Repo.Client.Set(context.Background(), "config", string(jsonConfig), 0).Err()
	if err != nil {
		w.Write([]byte("ошибка в чтении конфига"))
	}

}
