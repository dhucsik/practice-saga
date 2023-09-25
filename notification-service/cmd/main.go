package main

import (
	"context"
	"log"
	app2 "notification-service/app"
	"runtime"
)

func main() {
	ctx := context.Background()
	app := app2.InitApp(ctx)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}

	runtime.Goexit()
}
