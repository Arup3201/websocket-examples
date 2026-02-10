package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	HOST       = "0.0.0.0"
	PORT       = "8080"
	WRITE_WAIT = 10 * time.Second
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type Client struct {
	Id   string
	Conn *websocket.Conn
}

var clients = []*Client{}

func chat(w http.ResponseWriter, r *http.Request) {
	var err error
	var conn *websocket.Conn

	conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("[ERROR] upgrader upgrade: %s\n", err)
		return
	}

	id := uuid.NewString()
	clients = append(clients, &Client{
		Id:   id,
		Conn: conn,
	})

	fmt.Printf("[INFO] Client %s connected\n", id)

	defer func() {
		ind := slices.IndexFunc(clients, func(c *Client) bool {
			return c.Id == id
		})
		clients = append(clients[:ind], clients[ind+1:]...)

		err = conn.Close()
		if err != nil {
			fmt.Printf("[ERROR] connection close: %s\n", err)
			return
		}

		fmt.Printf("[INFO] Client %s connection closed\n", id)
	}()

	var mt, n int
	var data []byte
	var writer io.WriteCloser
	for {
		mt, data, err = conn.ReadMessage()
		if err != nil {
			fmt.Printf("[ERROR] connection read message: %s\n", err)
			break
		}

		if mt == websocket.TextMessage {
			for _, client := range clients {

				if client.Id == id {
					continue
				}

				err = client.Conn.SetWriteDeadline(time.Now().Add(WRITE_WAIT))
				if err != nil {
					fmt.Printf("[ERROR] connection set write deadline: %s\n", err)
					continue
				}

				writer, err = client.Conn.NextWriter(websocket.TextMessage)
				if err != nil {
					fmt.Printf("[ERROR] connection next writer: %s\n", err)
					continue
				}

				n, err = writer.Write(data)
				if err != nil {
					fmt.Printf("[ERROR] writer write: %s\n", err)
					continue
				}

				fmt.Printf("[INFO] written: %d bytes\n", n)

				if err = writer.Close(); err != nil {
					fmt.Printf("[ERROR] writer close: %s\n", err)
				}
			}
		}
	}
}

func main() {
	http.HandleFunc("/ws", chat)

	var port string
	port = os.Getenv("PORT")
	if port == "" {
		port = PORT
	}

	fmt.Printf("Server starting...\n")
	fmt.Printf("Host: %s\nPort: %s\n", HOST, port)

	err := http.ListenAndServe(fmt.Sprintf("%s:%s", HOST, port), nil)
	if err != nil {
		fmt.Printf("[ERROR] http listen and serve: %s\n", err)
	}
}
