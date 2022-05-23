package routers

import (
	"github.com/alirezakargar1380/agar.io-golang/app/api/adapter"
	"github.com/alirezakargar1380/agar.io-golang/app/api/endpoints"
	"github.com/alirezakargar1380/agar.io-golang/app/api/handlers"
	"github.com/alirezakargar1380/agar.io-golang/app/api/middlewares"
)

func ApiRouters() {
	/* TEST ROUTERS */
	Router.Handle("/m/t", adapter.Adapt(
		handlers.Test_main_handler(),
		middlewares.Function_num_2(),
		middlewares.Function_num_1(),
	)).Methods("POST")
	/* TEST ROUTERS */

	// Router.HandleFunc("/users/sigh_in", adapter.Adapt(

	// )).Methods("POST")

	/*	USER ROUTERS	*/
	Router.Handle("/users/sign_in", adapter.Adapt(
		handlers.Users_SignIn_Handler(),
		middlewares.Add_application_json_header(),
	)).Methods("POST")
	/*	USER ROUTERS	*/

	Router.Handle("/get/skins/{user_id}", adapter.Adapt(
		handlers.Gettt(),
		middlewares.Add_application_json_header(),
	)).Methods("GET")
	Router.HandleFunc("/add/skins", endpoints.AddSkinEndpoint).Methods("POST")
}
