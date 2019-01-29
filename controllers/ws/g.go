package ws

import (
	"flag"
	"log"
	"net/http"
)

var backend *Backend

func init() {
	flag.Parse()
	backend = newBackend()
	go backend.run()
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	serveWs(backend, w, r)
}

func serveWs(backend *Backend, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{backend: backend, conn: conn, send: make(chan []byte, 256)}
	client.backend.register <- client

	go client.writePump()
	go client.readPump()
}
