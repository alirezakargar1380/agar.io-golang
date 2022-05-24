package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alirezakargar1380/agar.io-golang/app/api/endpoints"
	"github.com/alirezakargar1380/agar.io-golang/app/api/routers"
	"github.com/alirezakargar1380/agar.io-golang/app/databases"
	"github.com/alirezakargar1380/agar.io-golang/app/socket"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
)

func main() {
	/*
		.env file
	*/
	err := godotenv.Load("config/config.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	databases.Initial_mongo_db()
	/*
		initialize redis db
	*/

	databases.Client = &databases.RedisDb{
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
	routers.ApiRouters()
	/*	setup Routers	*/

	srv := &http.Server{
		Handler:      routers.Router,
		Addr:         "localhost:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
