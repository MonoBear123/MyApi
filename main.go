package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/MonoBear123/MyApi/application"
)

func main() {
	app := application.New()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	err := app.Start(ctx)
	if err != nil {
		log.Fatal("server isn`t starting")
	}
	defer cancel()
}
