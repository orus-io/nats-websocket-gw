# NATS <-> websocket gateway

A websocket gateway for NATS, compatible with
backend for [websocket-nats](https://github.com/isobit/websocket-nats)

Features:

- TLS support
- Each NATS command is sent as a separate websocket message
- Provides a hook to change the CONNECT phase, allowing the http server to
  handle the connection itself (for example based on a cookie of the http request)

How does it differ from other nats-websocket servers ?

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
