package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/alirezakargar1380/agar.io-golang/app/databases"
	"github.com/alirezakargar1380/agar.io-golang/app/types/users_types"
	"github.com/alirezakargar1380/agar.io-golang/app/utils"
	"github.com/alirezakargar1380/agar.io-golang/app/validation"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Users_SignIn_endpoint(w http.ResponseWriter, r *http.Request) {
	body := &users_types.SignInRequest{}
	utils.ParseBody(r, body)
	validationErrors := validation.Sign_in_request_validation(body)
	if validationErrors != nil {
		w.Write(validationErrors)
		return
	}

	user := bson.D{
		{"username", body.Username},
		{"password", body.Password},
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := databases.UsersCollection.InsertOne(ctx, user)

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

func Get_All_Users_endpoint(w http.ResponseWriter, r *http.Request) {
	var allData []bson.M = []bson.M{}
	var result bson.M

	findOptions := options.Find()
	params := mux.Vars(r)
	user_id := params["page_number"]
	page_num, _ := strconv.ParseInt(user_id, 10, 64)
	pageNumber := page_num - 1
	findOptions.SetSkip(int64(pageNumber) * 10)
	findOptions.SetLimit(int64(pageNumber)*10 + 10)
	// findOptions.SetSort(bson.D{{}})

	showLoadedCursor, err := databases.UsersCollection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		panic(err)
	}

	for showLoadedCursor.Next(context.TODO()) {
		result = bson.M{}
		err := showLoadedCursor.Decode(&result)
		if err != nil {
			panic(err)
		}
		allData = append(allData, result)
	}
	fmt.Println(len(allData))

	resData, _ := json.Marshal(allData)
	w.Write(resData)
}
