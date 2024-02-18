package application

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/MonoBear123/MyApi/back/model"
)

func (a *App) loadConfig() {

	currentConfig := model.Config{
		Plus:           2,
		Minus:          2,
		Division:       2,
		Multiplication: 2,
		Construction:   2,
		MaxGorutines:   10,
	}
	res, err := json.Marshal(currentConfig)
	if err != nil {
		log.Fatal("ошибка в маршалинге конфига")
	}
	time.Sleep(2 * time.Second)
	err = a.rdb.Set(context.Background(), "config", string(res), 0).Err()
	if err != nil {
		log.Fatal(" стандартные значение конфига не установились")
	}

}
