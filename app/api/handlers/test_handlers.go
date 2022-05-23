package handlers

import (
	"fmt"
	"net/http"
)

func Test_main_handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("3")
		w.Write([]byte("main handler"))
	})
}
