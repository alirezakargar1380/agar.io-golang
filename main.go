package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/alirezakargar1380/agar.io-golang/app/beads"
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

	// set default agar size
	socket.Agars[Id] = &socket.AgarDetail{
		X:         200,
		Y:         200,
		Radius:    60,
		Speed:     0,
		Max_Speed: 5,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &socket.Client{
		Client_id: Id,
		RoomID:    roomId,
		Hub:       hub,
		Conn:      ws,
		Send:      make(chan []byte, 256),
	}
	client.Hub.Register <- client
	log.Println("Client successfully connected...")

	go client.WritePump()
	go client.ReadPump()
}

func setupRoutes() {
	hub := socket.NewHub()
	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsEndpoint(hub, w, r)
	})
}

func main() {
	fmt.Println("hello im backEnd agario")
	room := "game"
	c1 := make(chan bool)
	beads := &beads.Beads{
		Beads: make(map[string]map[string]int),
	}

	if beads.Beads[room] == nil {
		beads.Beads[room] = make(map[string]int)
	}

	var wg sync.WaitGroup
	makeCoordinate := func(room string) {
		for x := 0; x < 10000; x++ {
			for y := 0; y < 3; y++ {
				sx := fmt.Sprintf("%v", x)
				sy := fmt.Sprintf("%v", y)
				beads.Set(room, sx+"_"+sy)
			}
		}
		wg.Done()
		// c1 <- true
	}

	checkCoordinate := func(room string) {
		for x := 0; x < 10000; x++ {
			for y := 0; y < 3; y++ {
				sx := fmt.Sprintf("%v", x)
				sy := fmt.Sprintf("%v", y)
				beads.Exist(room, sx+"_"+sy)
			}
		}
		wg.Done()
		c1 <- true
	}

	wg.Add(2)
	go makeCoordinate(room)
	go checkCoordinate(room)
	wg.Wait()

	select {
	case <-c1:
		fmt.Println(len(beads.Beads[room]))
	}

	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
