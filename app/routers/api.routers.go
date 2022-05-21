package routers

import (
	"net/http"
)

func ApiRouters() {
	Router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	}).Methods("POST")
}
