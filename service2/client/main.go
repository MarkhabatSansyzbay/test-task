package main

import (
	"log"
	"net/http"
	"net/rpc"

	"service2/web"

	"github.com/go-chi/chi"
)

const serverPort = ":1234"

func main() {
	c, err := rpc.Dial("tcp", serverPort)
	if err != nil {
		log.Fatalf("can't connect to server: %v", err)
	}

	handler := new(web.Handler)
	handler.RpcClient = c

	r := chi.NewRouter()
	handler.InitRoutes(r)

	if err := http.ListenAndServe("localhost:8080", r); err != nil {
		log.Fatal(err)
	}
}
