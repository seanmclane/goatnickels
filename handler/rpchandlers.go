package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/seanmclane/goatnickels/block"
	"github.com/seanmclane/goatnickels/rpc"
	"log"
	"net/http"
	"net/url"
	"time"
)

var clients = make(map[*websocket.Conn]clientConfig) // connected clients
type clientConfig struct {
	Connected bool
	Subs      []string
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(r.RemoteAddr)

	clients[conn] = clientConfig{
		Connected: true,
		Subs:      []string{"block", "transaction", "vote"},
	}
	fmt.Println(clients[conn].Connected)
	defer conn.Close()

	for {
		var msg rpc.JsonRpcMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, conn)
			break
		}
		//handle message types with appropriate functions
		switch msg.Method {
		case "subscribe":
			handleSubs(conn, msg)
		case "transaction":
			fmt.Println("transaction message received")
			var txion block.Transaction
			err := json.Unmarshal(msg.Params, &txion)
			if err != nil {
				log.Printf("error: %v", err)
				delete(clients, conn)
				break
			}
			if txion.AddTransaction() {
				//TODO: change this to a response message type, so the client knows whether it was successful
				rpc.BroadcastChannel <- msg
			}
		case "vote":
			fmt.Println("vote message received")
			var vote block.Vote
			err := json.Unmarshal(msg.Params, &vote)
			if err != nil {
				log.Printf("error: %v", err)
				delete(clients, conn)
				break
			}
			if vote.AddVote() {
				//TODO: change this to a response message type, so the client knows whether it was successful
				rpc.BroadcastChannel <- msg
			}
		default:
			rpc.BroadcastChannel <- msg
		}

	}
	fmt.Println("Client disconnected")
}

func BroadcastMessages() {
	for {
		msg := <-rpc.BroadcastChannel

		for client, config := range clients {
			for _, sub := range config.Subs {
				//filter out messages that the client is not subscribed to
				if sub == msg.Method {
					err := client.WriteJSON(msg)
					if err != nil {
						log.Println("error:", err)
						client.Close()
						delete(clients, client)
					}
				}
			}
		}
	}
}

func handleSubs(conn *websocket.Conn, msg rpc.JsonRpcMessage) {
	var subs []string
	var j []string
	err := json.Unmarshal(msg.Params, &j)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	for client, config := range clients {
		if client == conn {
			subs = config.Subs
			for _, p := range j {
				subs = append(subs, p)
			}
		}
	}
	clients[conn] = clientConfig{
		Connected: true,
		Subs:      subs,
	}
}

func ConnectToNodes() {
	time.Sleep(1 * time.Second)
	config := block.LoadConfig()

	for _, node := range config.Nodes {
		go connectToNode(node)
	}
}

func connectToNode(node string) {
	u := url.URL{Scheme: "ws", Host: node + ":3000", Path: "api/v1/ws"}
	fmt.Println("connecting to node:", u.String())
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("error:", err)
	}
	defer conn.Close()
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", message)
	}
}
