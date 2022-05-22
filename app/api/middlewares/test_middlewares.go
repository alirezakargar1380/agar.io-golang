package middlewares

import (
	"fmt"
	"net/http"

	"github.com/alirezakargar1380/agar.io-golang/app/api/adapter"
)

func Function_num_2() adapter.Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("2")
			h.ServeHTTP(w, r)
		})
	}
}

func Function_num_1() adapter.Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("1")
			h.ServeHTTP(w, r)
		})
	}
}
