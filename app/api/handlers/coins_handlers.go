package handlers

import (
	"net/http"

	"github.com/alirezakargar1380/agar.io-golang/app/api/endpoints"
)

func Add_coin_handler() http.Handler {
	return http.HandlerFunc(endpoints.Add_coin_endpoints)
}
