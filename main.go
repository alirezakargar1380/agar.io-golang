package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/alirezakargar1380/agar.io-golang/app/socket"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsEndpoint(hub *socket.Hub, w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	roomId := params.Get("d")
	var clientId string = params.Get("client_id")
	Id, err := strconv.ParseInt(clientId, 0, 64)
	if err != nil {
		fmt.Println(err)
		return
	}

	socket.Agars[Id] = &socket.AgarDetail{
		Size: 1,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &socket.Client{
		Client_id: Id,
		// Client_id: 2,
		RoomID: roomId,
		Hub:    hub,
		Conn:   ws,
		Send:   make(chan []byte, 256),
	}
	client.Hub.Register <- client
	log.Println("Client successfully connected...")

	go client.WritePump()
	go client.ReadPump()
}

// 6037 - 6576 - 4606 - 6198
// 8.5

func setupRoutes() {
	hub := socket.NewHub()
	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsEndpoint(hub, w, r)
	})
}

type St struct {
	Name string
	Age  float64
}

func main() {
	fmt.Println("hello im backEnd agario")
	var person map[string]St = make(map[string]St)
	person["1"] = St{
		Name: "ali",
		Age:  20,
	}
	person["2"] = St{
		Name: "reza",
		Age:  30,
	}
	fmt.Println(person["2"].Age)
	person["2"] = St{
		Name: person["2"].Name,
		Age:  person["2"].Age + 0.1,
	}
	fmt.Println(person["2"].Age)
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
