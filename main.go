package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/MonoBear123/MyApi/server/application"
)

func main() {
	app := application.New()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	err := app.Start(ctx)
	if err != nil {
		log.Fatal(err)

	}
	defer cancel()
}
