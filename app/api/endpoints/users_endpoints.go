package endpoints

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/alirezakargar1380/agar.io-golang/app/databases"
	"github.com/alirezakargar1380/agar.io-golang/app/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"net/http"
)

type signInRequest struct {
	Username string
	Password string
}

func Users_SignIn_endpoint(w http.ResponseWriter, r *http.Request) {
	// GET DATA
	body := &signInRequest{}
	utils.ParseBody(r, body)

	database := databases.Mongodb_client.Database("agario")
	var usersCollection *mongo.Collection = database.Collection("users")
	user := bson.D{
		{"username", body.Username},
		{"password", body.Password},
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := usersCollection.InsertOne(ctx, user)

	if err != nil {
		log.Fatal(err)
	}
	resData, _ := json.Marshal(res)

	w.Write(resData)
}

func Users_SignUp_endpoint(w http.ResponseWriter, r *http.Request) {

}

func Get_Users_endpoint(w http.ResponseWriter, r *http.Request) {

}
