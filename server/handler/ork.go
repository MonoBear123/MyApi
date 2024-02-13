package handler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/MonoBear123/MyApi/back/model"
	"github.com/MonoBear123/MyApi/back/repository/expression"
	"github.com/go-chi/chi"
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
	w.Write([]byte(body.Expression1))
	parsedExpression, err := expression.ParseExpressionToTree(body.Expression1)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}
	num := expression.EvaluateExpression(parsedExpression)
	out := model.EXpression{
		Expression:  body.Expression1,
		ExpressinID: rand.Uint64(),
		Tree:        parsedExpression,
		Num:         num,
	}

	err = o.Repo.Insert(r.Context(), out)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		w.Write([]byte("ошибка в базе данных"))
		return
	}
	res, err := json.Marshal(out.ExpressinID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		w.Write([]byte("ошибка не в базе данных"))
		return
	}
	w.Write(res)
	w.WriteHeader(http.StatusCreated)

}
func (o *Expression) GetExpressionByID(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "id")

	const decimal = 10
	const bitSize = 64
	ID, err := strconv.ParseUint(key, decimal, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res, err := o.Repo.ExpressionFind(r.Context(), ID)
	if err != nil {
		return
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		fmt.Println("error in marshing")
	}

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
	status:=o.Repo.AgentALLFind(r.Context())
	for _, onestatus:=range status{
		timeOld:=time.Since(onestatus.Time)
		if timeOld>time.Second*5{
			w.Write([]byte("Агент умер\n"))
		}else{
			w.Write([]byte("Агент жив\n"))
		}

	}
}
