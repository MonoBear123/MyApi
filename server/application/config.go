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
		Minus:          100,
		Division:       100,
		Multiplication: 100,
		Construction:   100,
		MaxGorutines:   10,
	}
	res, err := json.Marshal(currentConfig)
	if err != nil {
		log.Fatal("ошибка в маршалинге конфига")
	}
	time.Sleep(15 * time.Second)
	err = a.rdb.Set(context.Background(), "config", string(res), 0).Err()
	if err != nil {
		log.Fatal(" стандартные значение конфига не установились")
	}

}
