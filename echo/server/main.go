package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	HOST          = "0.0.0.0"
	PORT          = "80"
	PING_INTERVAL = 2 * time.Second
	PING_WAIT     = 50 * time.Second
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
		log.Printf("[ERROR] upgrader: %s\n", err)
		return
	}
	defer func() {
		log.Printf("[INFO] Closing connection\n")
		conn.Close()
	}()

	log.Printf("[INFO] Client connected\n")

	go func() {
		ticker := time.NewTicker(PING_INTERVAL)
		defer ticker.Stop()

		conn.SetPongHandler(func(appData string) error {
			// PONG!
			conn.SetReadDeadline(time.Now().Add(PING_WAIT))
			return nil
		})
		for {
			<-ticker.C

			// PING!
			err = conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Printf("[ERROR] message ping: %s\n", err)
				conn.WriteMessage(websocket.CloseMessage, nil)
				break
			}
		}
	}()

	// read/write loop
	for {
		mt, msg, err = conn.ReadMessage()
		if err != nil {
			if err.(*websocket.CloseError).Code != websocket.CloseGoingAway {
				log.Printf("[ERROR] connection read: %s\n", err)
			}
			conn.WriteMessage(websocket.CloseMessage, nil)
			break
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
		}
	}
}

func main() {
	http.HandleFunc("/ws", serve)

	fmt.Printf("Server starting...\n")
	fmt.Printf("Host: %s\nPort: %s\n", HOST, PORT)
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", HOST, PORT), nil)
	if err != nil {
		log.Fatalf("http listen and serve: %s", err)
	}
}
