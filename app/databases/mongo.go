package databases

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Mongodb_client *mongo.Client
var Mongo_Context context.Context

func Initial_mongo_db() {
	fmt.Println("initial mongodb")
	mongo_url := os.Getenv("mongo_db_url")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongo_url))
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
