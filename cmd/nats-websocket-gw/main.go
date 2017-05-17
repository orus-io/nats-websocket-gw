package main

import (
	"net/http"

	"github.com/orus-io/nats-websocket-gw"
)

func main() {
	gateway := gw.NewGateway(gw.Settings{
		NatsAddr: "localhost:4222",
	})
	http.HandleFunc("/nats", gateway.Handler)
	http.ListenAndServe("0.0.0.0:8910", nil)
}
