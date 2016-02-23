package main

import (
	"github.com/micro/event-web/handler"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-web"

	event "github.com/micro/event-srv/proto/event"
)

func main() {
	service := web.NewService(
		web.Name("go.micro.web.event"),
		web.Handler(handler.Router()),
	)

	service.Init()

	handler.Init(
		"templates",
		event.NewEventClient("go.micro.srv.event", client.DefaultClient),
	)

	service.Run()
}
