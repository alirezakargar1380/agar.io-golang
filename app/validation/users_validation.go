package validation

import (
	"encoding/json"
	"fmt"

	"github.com/alirezakargar1380/agar.io-golang/app/types/users_types"
	"github.com/gookit/validate"
)

func Sign_in_request_validation(body *users_types.SignInRequest) []byte {
	v := validate.New(body)
	v.AddRule("username", "minLen", 4)
	v.AddRule("username", "regex", `^[a-zA-Z0-9]+[\s]`)

	v.AddRule("password", "minLen", 8)

	if v.Validate() {
		return nil
	} else {
		fmt.Println(v.Errors)
		rrr, _ := json.Marshal(v.Errors)
		return rrr
	}
}
