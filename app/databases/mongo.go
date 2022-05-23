package databases

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Mongodb_client *mongo.Client
var Mongo_Context context.Context

func Initial_mongo_db() {
	fmt.Println("initial mongodb")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	Mongodb_client = client
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = Mongodb_client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// defer Mongodb_client.Disconnect(ctx)
	// Mongo_Context = ctx
	Mongodb_client = client
}
