package main

import (
	"github.com/gorilla/mux"
	"github.com/micro/event-web/handler"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-web"

	event "github.com/micro/event-srv/proto/event"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handler.Index)
	r.HandleFunc("/search", handler.Search)
	r.HandleFunc("/latest", handler.Latest)
	r.HandleFunc("/event/{id}", handler.Event)

	service := web.NewService(
		web.Name("go.micro.web.event"),
		web.Handler(r),
	)

	service.Init()
	handler.EventClient = event.NewEventClient("go.micro.srv.event", client.DefaultClient)
	service.Run()
}
