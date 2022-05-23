package handlers

import (
	"net/http"

	"github.com/alirezakargar1380/agar.io-golang/app/api/endpoints"
)

func Gettt() http.Handler {
	return http.HandlerFunc(endpoints.GetSkinsEndpoint)
}
