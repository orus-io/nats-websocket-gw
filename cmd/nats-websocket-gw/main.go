package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/orus-io/nats-websocket-gw"
)

func usage() {
	fmt.Printf(`Usage: %s [ --help ] [ --no-origin-check ]
`, os.Args[0])
}

func main() {
	settings := gw.Settings{
		NatsAddr: "localhost:4222",
	}

	if len(os.Args) > 2 {
		fmt.Print("Too many arguments\n\n")
		usage()
		os.Exit(1)
	} else if len(os.Args) == 2 {
		if os.Args[1] == "--help" {
			usage()
			return
		} else if os.Args[1] == "--no-origin-check" {
			settings.WSUpgrader = &websocket.Upgrader{
				ReadBufferSize:  1024,
				WriteBufferSize: 1024,
				CheckOrigin:     func(r *http.Request) bool { return true },
			}
		} else {
			fmt.Printf("Invalid args: %s\n\n", os.Args[1])
			usage()
			return
		}
	}

	gateway := gw.NewGateway(settings)
	http.HandleFunc("/nats", gateway.Handler)
	http.ListenAndServe("0.0.0.0:8910", nil)
}
