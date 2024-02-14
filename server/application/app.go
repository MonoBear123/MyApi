package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type App struct {
	router http.Handler
	rdb    *redis.Client
}

func New() *App {
	app := &App{
		rdb: redis.NewClient(&redis.Options{
			Addr: "redis:6379",
		}),
	}

	app.loadRoutes()
	return app
}
func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    "0.0.0.0:8041",
		Handler: a.router,
	}

	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("not made")
	}

	defer func() {
		if err = a.rdb.Close(); err != nil {
			fmt.Print("redis close")
		}
	}()

	ch := make(chan error, 1)
	fmt.Println("Starting server ")

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- err
		}
		close(ch)
	}()

	select {
	case <-ctx.Done():
		timeout, canel := context.WithTimeout(context.Background(), time.Second*1000)
		defer canel()
		return server.Shutdown(timeout)
	case err = <-ch:
		return err
	}

}
