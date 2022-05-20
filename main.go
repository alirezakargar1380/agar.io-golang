package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	grpc_routers "github.com/alirezakargar1380/agar.io-golang/app/grpc"
	"github.com/alirezakargar1380/agar.io-golang/app/socket"
	"github.com/alirezakargar1380/agar.io-golang/app/trigonometric_circle"
	"github.com/alirezakargar1380/agar.io-golang/test_log"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
)

var upgrader = websocket.Upgrader{
	// ReadBufferSize:  512,
	// WriteBufferSize: 512,
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
			X:         60,
			Y:         60,
			Radius:    50,
			Name:      "",
			Max_speed: 7,
			Speed:     0,
		})

	} else {
		socket.Agars[roomId][Id].Color = "0xdfff994"
		socket.Agars[roomId][Id].Agars = append(socket.Agars[roomId][Id].Agars, trigonometric_circle.AgarDe{
			Lock:      false,
			Id:        1,
			X:         300,
			Y:         2900,
			Radius:    59,
			Name:      "",
			Max_speed: 7,
			Speed:     0,
		})

	}
	// }

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	var color string = ""
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
		Loose:     false,
	}
	client.Hub.Register <- client
	log.Println("Client successfully connected...", clientId)

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

	// redis_db.Client = &redis_db.RedisDb{
	// 	Client: redis.NewClient(&redis.Options{
	// 		Addr:     "127.0.0.1:6379",
	// 		Password: "",
	// 		DB:       0,
	// 	}),
	// }

	go func() {
		fmt.Println("hello im backEnd agario")
		setupRoutes()
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	s := grpc_routers.Serverr{}
	test_log.RegisterUserServiceServer(grpcServer, &s)

	fmt.Println("server grpc")
	grpcServer.Serve(lis)
}
