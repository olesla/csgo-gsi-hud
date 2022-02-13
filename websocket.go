package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	// constantly listen for new messages
	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
	}
}

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan []byte
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			break
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in Pool")
			for client, _ := range pool.Clients {
				if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// check the origin of our connection to allow us to make requests
	// from our vue development server to here. for now, do no checking
	// and just allow any connection
	CheckOrigin: func(r *http.Request) bool { return true },
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return conn, nil
}

func serveWs(pool *Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}
