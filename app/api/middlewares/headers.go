package middlewares

import (
	"fmt"
	"net/http"

	"github.com/alirezakargar1380/agar.io-golang/app/api/adapter"
)

func Add_application_json_header() adapter.Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Add_application_json_header middleware")
			w.Header().Set("Content-Type", "application/json")
			h.ServeHTTP(w, r)
		})
	}
}
