package middlewares

import (
	"fmt"
	"net/http"

	"github.com/alirezakargar1380/agar.io-golang/app/api/adapter"
)

func CheckThisUserIsCreatedBefore() adapter.Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("CheckThisUserIsCreatedBefore")
			h.ServeHTTP(w, r)
		})
	}
}
