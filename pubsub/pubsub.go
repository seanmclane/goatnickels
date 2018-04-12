package pubsub

import (
	"github.com/gorilla/websocket"
	"github.com/seanmclane/goatnickels/rpc"
)

var DefaultHub = NewHub()

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan rpc.JsonRpcMessage
	Register   chan *Client
	Unregister chan *Client
}

type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Subs     []string
	Id       int
	Requests []rpc.JsonRpcMessage
	Send     chan rpc.JsonRpcMessage
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan rpc.JsonRpcMessage),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}

func (c *Client) ReadPump(f func(c *Client)) {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	f(c)

}

func (c *Client) WritePump(f func(c *Client)) {
	defer c.Conn.Close()

	f(c)

}
