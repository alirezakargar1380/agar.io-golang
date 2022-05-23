package databases

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongodb_client *mongo.Client

func Initial_mongo_db() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	mongodb_client = client
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = mongodb_client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer mongodb_client.Disconnect(ctx)

}
