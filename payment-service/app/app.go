package app

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"payment-service/http"
	"payment-service/publisher"
	"payment-service/repository"
	"payment-service/service"
	"payment-service/subscriber"
)

type App struct {
	httpServer *http.Server
	repo       repository.Repository
	pub        publisher.Publisher
	srv        service.Service
	sub        subscriber.Subscriber
}

func InitApp(ctx context.Context) *App {
	app := &App{}

	for _, init := range []func(ctx context.Context) error{
		app.initRepository,
		app.initPublisher,
		app.initService,
		app.initSubscriber,
		app.initHTTPServer,
	} {
		err := init(ctx)
		if err != nil {
			log.Fatal("init app - ", err.Error())
			return nil
		}
	}

	fmt.Println(app.pub, app.repo, app.sub, app.srv)

	return app
}

func (a *App) initHTTPServer(_ context.Context) error {
	a.httpServer = http.NewServer(a.srv)

	return nil
}

func (a *App) initRepository(_ context.Context) error {
	var err error
	a.repo, err = repository.New()
	if err != nil {
		return err
	}

	return nil
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
	a.srv = service.New(a.repo, a.pub)

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
	g.Go(a.httpServer.Start)
	return g.Wait()
}
