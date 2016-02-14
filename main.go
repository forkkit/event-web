package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gorilla/mux"
	event "github.com/micro/event-srv/proto/event"
	"github.com/micro/event-web/handler"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/registry"
	"github.com/pborman/uuid"
)

func serve(r http.Handler) {
	l, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	go http.Serve(l, r)

	host, p, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	port, err := strconv.Atoi(p)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(host) == 0 || host == "[::]" || host == "::" {
		host, _ = os.Hostname()
	}

	service := &registry.Service{
		Name:    "go.micro.web.event",
		Version: "latest",
		Nodes: []*registry.Node{&registry.Node{
			Id:      uuid.NewUUID().String(),
			Address: host,
			Port:    port,
		}},
	}

	if err := registry.Register(service); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Listening on %v\n", l.Addr().String())

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	fmt.Printf("Received signal %s\n", <-ch)

	if err := registry.Deregister(service); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := l.Close(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	cmd.Init()

	handler.EventClient = event.NewEventClient("go.micro.srv.event", client.DefaultClient)

	r := mux.NewRouter()
	r.HandleFunc("/", handler.Index)
	r.HandleFunc("/search", handler.Search)
	r.HandleFunc("/latest", handler.Latest)
	r.HandleFunc("/event/{id}", handler.Event)
	serve(r)
}
