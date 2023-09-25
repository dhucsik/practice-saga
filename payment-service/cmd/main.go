package main

import (
	"context"
	"log"
	app2 "payment-service/app"
)

func main() {
	ctx := context.Background()
	app := app2.InitApp(ctx)

	log.Fatal(app.Start())
}
