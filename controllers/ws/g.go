package ws

import (
	"flag"
	"log"
	"net/http"
)

var backend *Backend

func init() {
	flag.Parse()
	backend = NewBackend()
	go backend.Run()
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

func TestWSHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "templates/test_ws.html")
}
