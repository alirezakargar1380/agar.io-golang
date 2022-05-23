package endpoints

import (
	"net/http"
)

func Users_SignIn_endpoint(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("si in"))
}
