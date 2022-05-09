package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/alirezakargar1380/agar.io-golang/app/socket"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// ReadBufferSize:  512,
	// WriteBufferSize: 512,
}

func wsEndpoint(hub *socket.Hub, w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	roomId := params.Get("d")
	var clientId string = params.Get("client_id")
	// Id, err := strconv.ParseInt(clientId, 0, 64)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// set default agar size
	// if socket.Agars[roomId] == nil {
	// 	socket.Agars[roomId] = make(map[int64]*socket.AgarDetail)
	// }

	// if socket.Agars[roomId][Id] == nil || len(socket.Agars[roomId][Id].Agars) == 0 {
	// socket.Agars[roomId][Id] = &socket.AgarDetail{
	// 	Client_id: int(Id),
	// }
	// if Id == 1 {
	// 	socket.Agars[roomId][Id].Color = "0xdfbg004"
	// 	socket.Agars[roomId][Id].Agars = append(socket.Agars[roomId][Id].Agars, trigonometric_circle.AgarDe{
	// 		Lock:      false,
	// 		Id:        1,
	// 		X:         100,
	// 		Y:         2900,
	// 		Radius:    50,
	// 		Name:      "",
	// 		Max_speed: 7,
	// 		Speed:     0,
	// 	}, trigonometric_circle.AgarDe{
	// 		Lock:      true,
	// 		Id:        2,
	// 		X:         100,
	// 		Y:         2750,
	// 		Radius:    60,
	// 		Name:      "",
	// 		Max_speed: 7,
	// 		Speed:     0,
	// 	})

	// } else {
	// 	socket.Agars[roomId][Id].Color = "0xdfff994"
	// 	socket.Agars[roomId][Id].Agars = append(socket.Agars[roomId][Id].Agars, trigonometric_circle.AgarDe{
	// 		Lock:      false,
	// 		Id:        1,
	// 		X:         300,
	// 		Y:         2900,
	// 		Radius:    60,
	// 		Name:      "",
	// 		Max_speed: 7,
	// 		Speed:     0,
	// 	})

	// }
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

	// var color string = ""
	// if Id == 1 {
	// 	color = "0xdfbg004"
	// } else {
	// 	color = "0xdfff994"
	// }

	client := &socket.Client{
		// Client_id: Id,
		RoomID: roomId,
		Hub:    hub,
		Conn:   ws,
		Send:   make(chan []byte, 256),
		// Color:  color,
	}
	client.Hub.Register <- client
	log.Println("Client successfully connected...", clientId)

	// go client.WritePump()
	// go client.ReadPump()
}

func setupRoutes() {
	hub := socket.NewHub()
	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsEndpoint(hub, w, r)
	})
}

// var stars map[string]bool = make(map[string]bool)

func main() {
	fmt.Println("hello im backEnd agario")
	var redisClient *redis.Client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	var stars map[string]bool = make(map[string]bool)
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			key := fmt.Sprintf("%d_%d", x, y)
			stars[key] = true
		}
	}

	ParseStars, err := json.Marshal(stars)
	if err != nil {
		fmt.Println(err)
	}

	err = redisClient.Set("stars", ParseStars, 0).Err()
	if err != nil {
		fmt.Println(err)
	}

	val, err := redisClient.Get("stars").Result()
	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal([]byte(val), &stars)
	fmt.Println(stars["1_99"])
	delete(stars, "1_99")
	fmt.Println(stars["1_99"])

	pp, err := json.Marshal(stars)
	if err != nil {
		fmt.Println(err)
	}

	err = redisClient.Set("stars", pp, 0).Err()
	if err != nil {
		fmt.Println(err)
	}

	vv, err := redisClient.Get("stars").Result()
	if err != nil {
		fmt.Println(err)
	}

	json.Unmarshal([]byte(vv), &stars)
	fmt.Println(stars["1_99"])

	return
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
