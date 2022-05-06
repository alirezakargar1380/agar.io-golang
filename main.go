package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/alirezakargar1380/agar.io-golang/app/socket"
	"github.com/alirezakargar1380/agar.io-golang/app/trigonometric_circle"
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
	if socket.Agars[roomId] == nil {
		socket.Agars[roomId] = make(map[int64]*socket.AgarDetail)
	}

	// if socket.Agars[roomId][Id] == nil || len(socket.Agars[roomId][Id].Agars) == 0 {
	socket.Agars[roomId][Id] = &socket.AgarDetail{
		Client_id: int(Id),
	}
	if Id == 1 {
		socket.Agars[roomId][Id].Color = "0xdfbg004"
		socket.Agars[roomId][Id].Agars = append(socket.Agars[roomId][Id].Agars, trigonometric_circle.AgarDe{
			Lock:      false,
			Id:        1,
			X:         100,
			Y:         100,
			Radius:    50,
			Name:      "",
			Max_speed: 7,
			Speed:     0,
		})

		// trigonometric_circle.AgarDe{
		// 	Lock:      true,
		// 	Id:        2,
		// 	X:         300,
		// 	Y:         100,
		// 	Radius:    60,
		// 	Name:      "",
		// 	Max_speed: 7,
		// 	Speed:     0,
		// }
	} else {
		socket.Agars[roomId][Id].Color = "0xdfff994"
		socket.Agars[roomId][Id].Agars = append(socket.Agars[roomId][Id].Agars, trigonometric_circle.AgarDe{
			Lock:      false,
			Id:        1,
			X:         100,
			Y:         250,
			Radius:    60,
			Name:      "",
			Max_speed: 7,
			Speed:     0,
		})

	}
	// }

	// index := 1
	// socket.Agars[roomId][Id].Agars = append(socket.Agars[roomId][Id].Agars[:index], socket.Agars[roomId][Id].Agars[index+1:]...)
	// fmt.Println(socket.Agars[roomId][Id].Agars)
	// return

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	var color string = ""
	fmt.Println(Id, "...")
	if Id == 1 {
		color = "0xdfbg004"
	} else {
		color = "0xdfff994"
	}

	client := &socket.Client{
		Client_id: Id,
		RoomID:    roomId,
		Hub:       hub,
		Conn:      ws,
		Send:      make(chan []byte, 256),
		Color:     color,
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

	// for i := 400; i > 0; i-- {
	var Radius float64 = float64(450)
	speed := 7 - (Radius * 0.013)
	fmt.Println(math.Floor(speed*1000) / 1000)
	// }

	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
