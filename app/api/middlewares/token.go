package middlewares

import (
	"net/http"

	"github.com/alirezakargar1380/agar.io-golang/app/api/adapter"
)

func CheckTokenisSet() adapter.Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// if user was admin
			h.ServeHTTP(w, r)
		})
	}
}

func CheckTokenisValidate() adapter.Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// if user was admin
			h.ServeHTTP(w, r)
		})
	}
}
