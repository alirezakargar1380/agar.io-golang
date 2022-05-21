package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alirezakargar1380/agar.io-golang/app/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Book struct {
	Name string
}

type SoldSkin struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Buyed_userId primitive.ObjectID `bson:"buyed_userId,omitempty"`
}

type Skin struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name,omitempty"`
	Sold_skin SoldSkin
}

func AddSkinEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

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
	var skinCollection *mongo.Collection = database.Collection("skins")
	user := bson.D{
		{"name", "test_skin"},
	}

	res, err := skinCollection.InsertOne(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
	resData, _ := json.Marshal(res)

	w.Write(resData)
}

func GetSkinsEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user_id := params["user_id"]
	id, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	body := &Book{}
	utils.ParseBody(r, body)
	fmt.Println(body)

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

	showLoadedCursor, err := skinCollection.Aggregate(context.TODO(), []bson.M{
		bson.M{"$lookup": bson.M{
			"from":         "sold_skins",
			"localField":   "_id",
			"foreignField": "skin_id",
			"as":           "sold_skin",
		}},
		bson.M{"$unwind": bson.M{
			"path": "$sold_skin",
		}},
		bson.M{"$match": bson.M{
			// "sold_skin.buyed_userId": bson.M{"$ne": id},
			"sold_skin.buyed_userId": bson.M{
				"$in": []primitive.ObjectID{id},
			},
		}},
		// bson.M{
		// 	"$group": bson.M{
		// 		"_id": "$_id",
		// 	}},
	})
	if err != nil {
		panic(err)
	}

	// lookupStage := bson.D{
	// 	{
	// 		"$lookup", bson.D{{"from", "sold_skins"}, {"localField", "_id"}, {"foreignField", "skin_id"}, {"as", "sold_skin"}},
	// 	},
	// }
	// unwindStage := bson.D{
	// 	{
	// 		"$unwind", bson.D{{"path", "$sold_skin"}, {"preserveNullAndEmptyArrays", true}},
	// 	},
	// }

	// showLoadedCursor, err := skinCollection.Aggregate(ctx, mongo.Pipeline{lookupStage, unwindStage,
	// 	bson.D{{"$match", bson.D{
	// 		{
	// 			"sold_skin.buyed_userId", bson.D{
	// 				{
	// 					"$eq", id,
	// 				},
	// 			},
	// 		},
	// 	}}},
	// })

	if err != nil {
		panic(err)
	}
	// var showsLoaded []bson.M
	var showsLoaded []Skin
	if err = showLoadedCursor.All(ctx, &showsLoaded); err != nil {
		panic(err)
	}
	// fmt.Println(showsLoaded)

	data, _ := json.Marshal(showsLoaded)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
