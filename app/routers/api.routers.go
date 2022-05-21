package routers

import (
	"fmt"
	"net/http"

	"github.com/alirezakargar1380/agar.io-golang/app/utils"
)

type Book struct {
	Name string
}

func ApiRouters() {
	Router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		body := &Book{}
		utils.ParseBody(r, body)
		fmt.Println(body)
		w.Write([]byte("test"))
	}).Methods("POST")
}
