package routers

import (
	"fmt"
	"net/http"

	"github.com/alirezakargar1380/agar.io-golang/app/api/adapter"
	"github.com/alirezakargar1380/agar.io-golang/app/api/middlewares"
	"github.com/alirezakargar1380/agar.io-golang/app/endpoints"
)

func MainHAN() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("3")
		w.Write([]byte("main handler"))
	})
}

func ApiRouters() {
	/* TEST ROUTERS */
	Router.Handle("/m/t", adapter.Adapt(
		MainHAN(),
		middlewares.Function_num_2(),
		middlewares.Function_num_1(),
	)).Methods("POST")
	/* TEST ROUTERS */

	Router.HandleFunc("/get/skins/{user_id}", endpoints.GetSkinsEndpoint).Methods("GET")
	Router.HandleFunc("/add/skins", endpoints.AddSkinEndpoint).Methods("POST")
}
