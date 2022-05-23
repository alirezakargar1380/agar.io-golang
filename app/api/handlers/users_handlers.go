package handlers

import (
	"net/http"

	"github.com/alirezakargar1380/agar.io-golang/app/api/endpoints"
)

func Users_SignIn_Handler() http.Handler {
	return http.HandlerFunc(endpoints.Users_SignIn_endpoint)
}
