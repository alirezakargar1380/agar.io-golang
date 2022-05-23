package routers

import (
	"github.com/alirezakargar1380/agar.io-golang/app/api/endpoints"
)

// func SocketRouters(w http.ResponseWriter, r *http.Request) {
// 	endpoints.SocketEndpoint(w, r)
// }

func SocketRouters() {
	Router.HandleFunc("/wss", endpoints.SocketEndpoint).Methods("GET")
}
