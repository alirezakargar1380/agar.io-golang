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
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func main() {
	expirationTime := time.Now().Add(1 * time.Minute)
	claims := &Claims{
		Username: "alireza",
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println(token)
	var jwtKey = []byte("my_secret_key")
	tokenString, errorsd := token.SignedString(jwtKey)
	if errorsd != nil {
		fmt.Println(errorsd)
	}
	fmt.Println(tokenString)

	claims = &Claims{}
	tknStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFsaXJlemEiLCJleHAiOjE2NTM0MTAzNTd9.Tq5pEc4DVKaoGqsfWeapvdYWpE_E0Su8dLgjf2ulrg4"
	tkn, errrrr := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if errrrr != nil {
		if errrrr == jwt.ErrSignatureInvalid {
			fmt.Println(http.StatusUnauthorized, "StatusUnauthorized")
			return
		}
		fmt.Println(http.StatusBadRequest, "StatusBadRequest")
		return
	}
	if !tkn.Valid {
		fmt.Println(http.StatusUnauthorized, "StatusUnauthorized")
		return
	}
	fmt.Println(claims.Username)

	/*   .env file   */
	err := godotenv.Load("config/config.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	/*	mongo_db    */
	databases.Initial_mongo_db()

	/*  initialize redis db	  */
	databases.Client = &databases.RedisDb{
		Client: redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       0,
		}),
	}

	/*	initialize Socket Hub	*/
	endpoints.Hub = socket.NewHub()
	go endpoints.Hub.Run()

	/*	setup Routers	*/
	fmt.Println("hello im backEnd agario")
	routers.SocketRouters()
	routers.ApiRouters()

	srv := &http.Server{
		Handler:      routers.Router,
		Addr:         "localhost:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
