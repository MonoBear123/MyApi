package application

import (
	"context"
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

	time.Sleep(15 * time.Second)
	_, err := a.rdb.Set(context.Background(), "config", currentConfig, 0).Result()
	if err != nil {
		log.Fatal(" стандартные значение конфига не установились")
	}

}
