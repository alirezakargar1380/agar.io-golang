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

	/*	USER ROUTERS	*/
	Router.Handle("/users/sign_in", adapter.Adapt(
		handlers.Users_SignIn_Handler(),
		middlewares.Add_application_json_header(),
	)).Methods("POST")

	Router.Handle("/users/sign_up", adapter.Adapt(
		handlers.Users_SignUp_Handler(),
		middlewares.Add_application_json_header(),
	)).Methods("POST")

	Router.Handle("/users/get_users/{page_number}", adapter.Adapt(
		handlers.Get_Users_Handler(),
		middlewares.Add_application_json_header(),
		middlewares.CheckIsAdmin(),
		middlewares.CheckTokenisValidate(),
	)).Methods("GET")
	/*	USER ROUTERS	*/

	/*	SKIN ROUTERS	*/
	Router.Handle("/get/skins/{user_id}", adapter.Adapt(
		handlers.Gettt(),
		middlewares.Add_application_json_header(),
	)).Methods("GET")
	Router.HandleFunc("/add/skins", endpoints.AddSkinEndpoint).Methods("POST")
	/*	SKIN ROUTERS	*/
}
