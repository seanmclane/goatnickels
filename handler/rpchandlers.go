package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/seanmclane/goatnickels/block"
	"github.com/seanmclane/goatnickels/pubsub"
	"github.com/seanmclane/goatnickels/rpc"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

/*
var clients = make(map[*websocket.Conn]clientConfig) // connected clients

type clientConfig struct {
	Connected bool
	Subs      []string
	Id        int
	Requests  []rpc.JsonRpcMessage
}
*/

type getBlockRequest struct {
	Index int
}

func HandleConnections(hub *pubsub.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(r.RemoteAddr)

	client := &pubsub.Client{Hub: hub, Conn: conn, Send: make(chan rpc.JsonRpcMessage)}
	client.Hub.Register <- client

	client.Subs = []string{"block", "transaction", "vote", "getblock"}
	client.Id = 1

	go client.ReadPump(ReadPumpFunc())
	go client.WritePump(WritePumpFunc())
}

func WritePumpFunc() func(c *pubsub.Client) {
	return func(c *pubsub.Client) {
		for {
			msg := <-c.Send
			for _, sub := range c.Subs {
				//filter out messages that the client is not subscribed to
				if sub == msg.Method {
					err := c.Conn.WriteJSON(msg)
					if err != nil {
						log.Println("error:", err)
					}
				}
			}
			//send result messages
			if msg.Id > 0 {
				err := c.Conn.WriteJSON(msg)
				if err != nil {
					log.Println("error:", err)
				}
			}
		}
	}
}

func ReadPumpFunc() func(c *pubsub.Client) {
	return func(c *pubsub.Client) {
		for {
			var msg rpc.JsonRpcMessage
			err := c.Conn.ReadJSON(&msg)
			if err != nil {
				log.Printf("error: %v", err)
				break
			}
			fmt.Println(msg.Method)
			//handle message types with appropriate functions
			switch msg.Method {
			case "subscribe":
				handleSubs(c, msg)
			case "transaction":
				block.TransactionChannel <- msg
			case "vote":
				block.VoteChannel <- msg
			case "getblock":
				handleGetBlock(c, msg)
			case "block":
				block.BlockChannel <- msg
			case "":
				handleResult(c, msg)
			default:
			}
		}
	}
}

func handleSubs(c *pubsub.Client, msg rpc.JsonRpcMessage) {
	var subs []string
	var j []string
	err := json.Unmarshal(msg.Params, &j)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	subs = c.Subs
	for _, p := range j {
		subs = append(subs, p)
	}
	c.Subs = subs
	c.Id = msg.Id

	out := []byte(`{"success": true}`)
	res := rpc.BuildResponse(msg.Id, out, nil)
	//send directly to client rather than putting on the broadcast channel
	c.Send <- res
}

func handleGetBlock(c *pubsub.Client, msg rpc.JsonRpcMessage) {
	fmt.Println("getblock request message received")
	var req getBlockRequest
	err := json.Unmarshal(msg.Params, &req)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	block := block.ReadBlockFromLocalStorage(strconv.Itoa(req.Index))
	bytes, err := json.Marshal(block)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	res := rpc.BuildResponse(msg.Id, bytes, nil)
	//send directly to client rather than putting on the broadcast channel
	c.Send <- res
}

func handleResult(c *pubsub.Client, msg rpc.JsonRpcMessage) {
	for _, req := range c.Requests {
		fmt.Println(req.Id, msg.Id)
		if msg.Id == req.Id {
			switch req.Method {
			case "getblock":
				var block block.Block
				err := json.Unmarshal(msg.Result, &block)
				if err != nil {
					log.Printf("error: %v", err)
					return
				}
				//TODO: save to local storage
				fmt.Println(block)
			default:
				fmt.Println("no method found")
			}
		}
	}
}

/*
func handleConnectionsFor(conn *websocket.Conn) {
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
			handleTransaction(conn, msg)
		case "vote":
			handleVote(conn, msg)
		case "getblock":
			handleGetBlock(conn, msg)
		case "":
			handleResult(conn, msg)
		default:
		}
	}
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
			if msg.Id > 0 {
				clients[client] = clientConfig{
					Connected: config.Connected,
					Subs:      config.Subs,
					Id:        msg.Id,
					Requests:  append(config.Requests, msg),
				}
				fmt.Println(clients[client].Requests)
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
		Id:        msg.Id,
	}
	out := []byte(`{"success": true}`)
	res := rpc.BuildResponse(msg.Id, out, nil)
	//send directly to client rather than putting on the broadcast channel
	conn.WriteJSON(res)
}

func handleTransaction(conn *websocket.Conn, msg rpc.JsonRpcMessage) {
	fmt.Println("transaction message received")
	var txion block.Transaction
	err := json.Unmarshal(msg.Params, &txion)
	if err != nil {
		log.Printf("error: %v", err)
		delete(clients, conn)
		return
	}
	if txion.AddTransaction() {
		//TODO: change this to a response message type, so the client knows whether it was successful
		rpc.BroadcastChannel <- msg
	}
}

func handleVote(conn *websocket.Conn, msg rpc.JsonRpcMessage) {
	fmt.Println("vote message received")
	var vote block.Vote
	err := json.Unmarshal(msg.Params, &vote)
	if err != nil {
		log.Printf("error: %v", err)
		delete(clients, conn)
		return
	}
	if vote.AddVote() {
		//TODO: change this to a response message type, so the client knows whether it was successful
		rpc.BroadcastChannel <- msg
	}
}

func handleGetBlock(conn *websocket.Conn, msg rpc.JsonRpcMessage) {
	fmt.Println("getblock request message received")
	var req getBlockRequest
	err := json.Unmarshal(msg.Params, &req)
	if err != nil {
		log.Printf("error: %v", err)
		delete(clients, conn)
		return
	}
	block := block.ReadBlockFromLocalStorage(strconv.Itoa(req.Index))
	bytes, err := json.Marshal(block)
	if err != nil {
		log.Printf("error: %v", err)
		delete(clients, conn)
		return
	}
	res := rpc.BuildResponse(msg.Id, bytes, nil)
	//send directly to client rather than putting on the broadcast channel
	conn.WriteJSON(res)
}

func handleResult(conn *websocket.Conn, msg rpc.JsonRpcMessage) {
	for _, req := range clients[conn].Requests {
		fmt.Println(req.Id, msg.Id)
		if msg.Id == req.Id {
			switch req.Method {
			case "getblock":
				var block block.Block
				err := json.Unmarshal(msg.Result, &block)
				if err != nil {
					log.Printf("error: %v", err)
					delete(clients, conn)
					return
				}
				fmt.Println(block)
			default:
				fmt.Println("no method found")
			}
		}
	}
}
*/

func ConnectToNodes(h *pubsub.Hub) {
	//sleep in case servers are not up yet, mostly for testing
	time.Sleep(1 * time.Second)
	config := block.LoadConfig()

	for _, node := range config.Nodes {
		go connectToNode(h, node)
	}
}

func connectToNode(hub *pubsub.Hub, node string) {
	u := url.URL{Scheme: "ws", Host: node + ":3000", Path: "api/v1/ws"}
	fmt.Println("connecting to node:", u.String())
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("error:", err)
	}

	client := &pubsub.Client{Hub: hub, Conn: conn, Send: make(chan rpc.JsonRpcMessage)}
	client.Hub.Register <- client

	client.Subs = []string{"block", "transaction", "vote", "getblock"}
	client.Id = 1

	go client.ReadPump(ReadPumpFunc())
	go client.WritePump(WritePumpFunc())
}
