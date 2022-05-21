package routers

import (
	"github.com/alirezakargar1380/agar.io-golang/app/endpoints"
)

func ApiRouters() {
	Router.HandleFunc("/get/skins/{user_id}", endpoints.GetSkinsEndpoint).Methods("GET")
	Router.HandleFunc("/add/skins", endpoints.AddSkinEndpoint).Methods("POST")
}
