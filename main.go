package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	redis_db "github.com/alirezakargar1380/agar.io-golang/app/service"
	"github.com/alirezakargar1380/agar.io-golang/app/socket"
	"github.com/alirezakargar1380/agar.io-golang/app/trigonometric_circle"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

type Sold_skin struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Buyed_userId primitive.ObjectID `bson:"buyed_userId,omitempty"`
	Skin_id      primitive.ObjectID `bson:"skin_id,omitempty"`
}

type Skin struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name primitive.ObjectID `bson:"name,omitempty"`
}

func main() {
	fmt.Println("hello im backEnd agario")

	redis_db.Client = &redis_db.RedisDb{
		Client: redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       0,
		}),
	}

	// pong, err := redis_db.Client.Client.Ping().Result()
	// if err == nil {
	// 	fmt.Println(pong)
	// } else {
	// 	panic(err)
	// }

	// TEST MONGO DB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	database := client.Database("agario")
	skinCollection := database.Collection("skins")
	// database2 := client.Database("quickstart")
	// episodesCollection := database2.Collection("episodes")
	// podcastCollection := database2.Collection("podcasts")

	/* Insert document */
	// user_id, _ := primitive.ObjectIDFromHex("628755681570972e5901c8e6")
	// skin_id, _ := primitive.ObjectIDFromHex("62876b848cf65ce052830d4e")
	// user := bson.D{
	// 	{"buyed_userId", user_id},
	// 	{"skin_id", skin_id},
	// }
	// collection.InsertOne(ctx, user)
	/* End Insert document */

	/* Test Aggregate */
	// id, _ := primitive.ObjectIDFromHex("5e3b37e51c9d4400004117e6")
	// lookupStage := bson.D{{"$lookup", bson.D{{"from", "podcasts"}, {"localField", "podcast"}, {"foreignField", "_id"}, {"as", "podcast"}}}}
	// unwindStage := bson.D{{"$unwind", bson.D{{"path", "$podcast"}, {"preserveNullAndEmptyArrays", false}}}}

	// showLoadedCursor, err := episodesCollection.Aggregate(ctx, mongo.Pipeline{lookupStage, unwindStage})
	// if err != nil {
	// 	panic(err)
	// }
	// var showsLoaded []bson.M
	// if err = showLoadedCursor.All(ctx, &showsLoaded); err != nil {
	// 	panic(err)
	// }
	// fmt.Println(showsLoaded)
	/* End Test Aggregate */

	// -------------------------------------------------------------------------------------------------------------

	/* Test Aggregate */
	// lookupStage := bson.D{{"$lookup", bson.D{{"from", "episodes"}, {"localField", "_id"}, {"foreignField", "podcast"}, {"as", "episode"}}}}
	// unwindStage := bson.D{{"$unwind", bson.D{{"path", "$episode"}, {"preserveNullAndEmptyArrays", false}}}}

	// showLoadedCursor, err := podcastCollection.Aggregate(ctx, mongo.Pipeline{lookupStage, unwindStage})
	// if err != nil {
	// 	panic(err)
	// }
	// var showsLoaded []bson.M
	// if err = showLoadedCursor.All(ctx, &showsLoaded); err != nil {
	// 	panic(err)
	// }
	// fmt.Println(showsLoaded)
	/* End Test Aggregate */

	/* Test Aggregate */
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "sold_skins"}, {"localField", "_id"}, {"foreignField", "skin_id"}, {"as", "sold_skin"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$sold_skin"}, {"preserveNullAndEmptyArrays", false}}}}

	showLoadedCursor, err := skinCollection.Aggregate(ctx, mongo.Pipeline{lookupStage, unwindStage})
	if err != nil {
		panic(err)
	}
	var showsLoaded []bson.M
	if err = showLoadedCursor.All(ctx, &showsLoaded); err != nil {
		panic(err)
	}
	fmt.Println(showsLoaded)
	/* End Test Aggregate */

	// END - TEST MONGO DB

	return
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
