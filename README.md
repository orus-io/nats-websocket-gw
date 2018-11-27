# NATS <-> websocket gateway

A websocket gateway for NATS, usable as a backend for
[elm-nats](https://github.com/orus-io/elm-nats) and
[websocket-nats](https://github.com/isobit/websocket-nats).

Features:

- TLS support
- Each NATS command is sent as a separate websocket message
- Provides a hook to change the CONNECT phase, allowing the http server to
  handle the connection itself (for example based on a cookie of the http request)
- Easily embeddable in a bigger http server
- Supports both text (default) and binary (by adding '?mode=binary' to the url) messages

## Basic usage

Fetch the source:

```bash
go get -u github.com/orus-io/nats-websocket-gw
```

Install and run the default binary

```bash
go install github.com/orus-io/nats-websocket-gw/cmd/nats-websocket-gw
nats-websocket-gw --no-origin-check
```

and/or integrate it in your http server:

```go
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
```

## How does it differ from other nats-websocket servers ?

- [Rest to NATS Proxy](https://github.com/sohlich/nats-proxy) provides a websocket
  based implementation. The approach is pretty different though, as the websockets
  do not transport the whole NATS protocol, but only the data of the messages
  (websockets are opened for each subscribed subject).

- [ws-tcp-relay](https://github.com/isobit/ws-tcp-relay), which the proposed
  backend for [websocket-nats](https://github.com/isobit/websocket-nats) is a generic
  ws-tcp gateway. As such, it:
  - does not allow for TLS connections to NATS because the TLS negociation
    cannot be done immediately (the NATS protocol has a clear text 'INFO' exchange
    before TLS handshake)
  - send websocket messages that may contain several NATS commands.
