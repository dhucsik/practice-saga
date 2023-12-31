package app

import (
	"context"
	"golang.org/x/sync/errgroup"
	"log"
	"notification-service/publisher"
	"notification-service/service"
	"notification-service/subscriber"
)

type App struct {
	pub publisher.Publisher
	srv service.Service
	sub subscriber.Subscriber
}

func InitApp(ctx context.Context) *App {
	app := &App{}

	for _, init := range []func(ctx context.Context) error{
		app.initPublisher,
		app.initService,
		app.initSubscriber,
	} {
		err := init(ctx)
		if err != nil {
			log.Fatal("init app - ", err.Error())
			return nil
		}
	}

	return app
}

func (a *App) initPublisher(_ context.Context) error {
	var err error
	a.pub, err = publisher.New()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initService(_ context.Context) error {
	a.srv = service.New(a.pub)
	return nil
}

func (a *App) initSubscriber(_ context.Context) error {
	var err error
	a.sub, err = subscriber.New(a.srv)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) Start() error {
	g := errgroup.Group{}

	g.Go(a.sub.Start)

	return g.Wait()
}
