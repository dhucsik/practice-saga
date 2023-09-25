package app

import (
	"context"
	"golang.org/x/sync/errgroup"
	"log"
	"order-service/http"
	"order-service/publisher"
	"order-service/repository"
	"order-service/service"
	"order-service/subscriber"
	"order-service/worker"
)

type App struct {
	httpServer *http.Server
	pub        publisher.Publisher
	repo       repository.Repository
	srv        service.Service
	sub        subscriber.Subscriber
	worker     worker.PusherWorker
}

func InitApp(ctx context.Context) *App {
	app := &App{}

	for _, init := range []func(ctx context.Context) error{
		app.initRepository,
		app.initPublisher,
		app.initService,
		app.initSubscriber,
		app.initHTTPServer,
		app.initWorker,
	} {
		err := init(ctx)
		if err != nil {
			log.Fatal("init app - ", err.Error())
			return nil
		}
	}

	return app
}

func (a *App) initWorker(_ context.Context) error {
	a.worker = worker.New(a.srv)
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

func (a *App) initHTTPServer(_ context.Context) error {
	a.httpServer = http.NewServer(a.srv)

	return nil
}

func (a *App) Start() error {
	g := errgroup.Group{}

	g.Go(a.worker.Start)
	g.Go(a.sub.Start)
	g.Go(a.httpServer.Start)

	return g.Wait()
}
