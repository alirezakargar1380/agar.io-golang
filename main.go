package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alirezakargar1380/agar.io-golang/app/endpoints"
	"github.com/alirezakargar1380/agar.io-golang/app/routers"
	redis_db "github.com/alirezakargar1380/agar.io-golang/app/service"
	"github.com/alirezakargar1380/agar.io-golang/app/socket"
	"github.com/go-redis/redis"
)

func main() {
	/*
		initialize redis db
	*/

	redis_db.Client = &redis_db.RedisDb{
		Client: redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       0,
		}),
	}

	/*
		initialize redis db
	*/

	/*	initialize Socket Hub	*/
	endpoints.Hub = socket.NewHub()
	go endpoints.Hub.Run()
	/*	initialize Socket Hub	*/

	/*	setup Routers	*/
	fmt.Println("hello im backEnd agario")
	routers.SocketRouters()
	/*	setup Routers	*/

	srv := &http.Server{
		Handler:      routers.Router,
		Addr:         "localhost:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
