package handler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/MonoBear123/MyApi/model"
	"github.com/MonoBear123/MyApi/repository/expression"
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

		w.WriteHeader(http.StatusBadRequest)
		return
	}
	now := time.Now()
	out := model.EXpression{
		Expression:  body.Expression1,
		ExpressinID: rand.Uint64(),
		CreatedEX:   &now,
	}

	err := o.Repo.Insert(r.Context(), out)
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
