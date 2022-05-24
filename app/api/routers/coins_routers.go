package routers

import (
	"github.com/alirezakargar1380/agar.io-golang/app/api/adapter"
	"github.com/alirezakargar1380/agar.io-golang/app/api/handlers"
	"github.com/alirezakargar1380/agar.io-golang/app/api/middlewares"
)

func Coins() {
	Router.Handle("/coins/add", adapter.Adapt(
		handlers.Add_coin_handler(),
		middlewares.Add_application_json_header(),
		middlewares.CheckIsAdmin(),
	)).Methods("POST")

	// Router.HandleFunc("/coins/get", endpoints.GetCoinsEndpoint).Methods("GET")

	// Router.HandleFunc("/coins/get/{user_id}", endpoints.GetCoinEndpoint).Methods("GET")

	// Router.HandleFunc("/coins/update/{user_id}", endpoints.UpdateCoinEndpoint).Methods("PUT")

	// Router.HandleFunc("/coins/delete/{user_id}", endpoints.DeleteCoinEndpoint).Methods("DELETE")
}
