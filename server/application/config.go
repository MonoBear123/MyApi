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
		Plus:           100,
		Minus:          100,
		Division:       100,
		Multiplication: 100,
		Construction:   100,
	}
	jsonConfig, err := json.Marshal(currentConfig)
	if err != nil {
		return
	}
	time.Sleep(15*time.Second)
	_, err = a.rdb.Set(context.Background(), "config", string(jsonConfig), 0).Result()
	if err != nil {
		log.Fatal(" стандартные значение конфига не установились")
	}

}
