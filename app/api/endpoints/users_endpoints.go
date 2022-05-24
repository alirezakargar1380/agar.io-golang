package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/alirezakargar1380/agar.io-golang/app/databases"
	"github.com/alirezakargar1380/agar.io-golang/app/types/users_types"
	"github.com/alirezakargar1380/agar.io-golang/app/utils"
	"github.com/gookit/validate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"net/http"
)

type signInRequest struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

func Users_SignIn_endpoint(w http.ResponseWriter, r *http.Request) {
	// GET DATA
	body := &signInRequest{}

	utils.ParseBody(r, body)
	fmt.Println(body.Username)

	v := validate.New(body)

	v.AddRule("username", "minLen", 7)
	v.AddRule("password", "minLen", 7)

	if v.Validate() {
		// validate ok
		fmt.Println("hello world")
	} else {
		rrr, _ := json.Marshal(v.Errors)
		w.Write([]byte(rrr))
		fmt.Println(v.Errors) // all error messages
	}
	return

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
	database := databases.Mongodb_client.Database("agario")
	var usersCollection *mongo.Collection = database.Collection("users")

	body := &users_types.SignUpRequerst{}
	utils.ParseBody(r, body)
	fmt.Println(body.Username)

	// err := usersCollection.FindOne(context.TODO(), bson.D{
	// 	{"username", "alireza"},
	// })

	var result bson.M
	err := usersCollection.FindOne(context.TODO(), bson.M{
		"username": "",
	}).Decode(&result)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}

func Get_Users_endpoint(w http.ResponseWriter, r *http.Request) {

}
