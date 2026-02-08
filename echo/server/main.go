package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	HOST = "localhost"
	PORT = 8080
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HTTP will automatically create multiple goroutines for multiple requests
// For every request websocket connection will be created and read/write loop
// will be started to read message from connection and write back to it
func serve(w http.ResponseWriter, r *http.Request) {
	var conn *websocket.Conn
	var err error
	var n int
	var msg []byte
	var wr io.WriteCloser
	var mt int

	conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ERROR] upgrader: %s", err)
		return
	}
	defer func() {
		log.Printf("[INFO] Closing connection\n")
		conn.Close()
	}()

	// read/write loop
	for {
		mt, msg, err = conn.ReadMessage()
		if err != nil {
			log.Printf("[ERROR] connection read: %s\n", err)
			break
		}

		switch mt {
		case websocket.BinaryMessage:
			log.Printf("[INFO] binary message received...\n")
		case websocket.TextMessage:
			log.Printf("[INFO] text message recieved: %s", msg)
		}

		wr, err = conn.NextWriter(mt)
		if err != nil {
			log.Printf("[ERROR] connection next writer: %s\n", err)
			continue
		}
		n, err = wr.Write(msg)
		if err != nil {
			log.Printf("[ERROR] connection write: %s\n", err)
			continue
		}

		log.Printf("[INFO] write bytes: %d\n", n)

		if err = wr.Close(); err != nil {
			log.Printf("[ERROR] connection write close: %s\n", err)
			continue
		}
	}
}

func main() {
	http.HandleFunc("/ws", serve)

	fmt.Printf("Server starting...\n")
	fmt.Printf("Host: %s\nPort: %d\n", HOST, PORT)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", HOST, PORT), nil)
	if err != nil {
		log.Fatalf("http listen and serve: %s", err)
	}
}
