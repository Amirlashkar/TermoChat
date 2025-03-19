package routers

import (
	"github.com/gorilla/mux"
)

// Main function to provide last router
func ProvideRouter() *mux.Router {
    router := mux.NewRouter()

    // ------------------ AUTH ------------------
    authRout := router.PathPrefix("/auth").Subrouter()
    authRout.Use(formMiddleWare)

	authRout.HandleFunc("/signup",   user_signup).Methods("POST")
	authRout.HandleFunc("/login",    user_login).Methods("POST")
	authRout.HandleFunc("/uexist",   user_existance).Methods("POST")
	// ------------------------------------------

	// ------------------ USERS ROUTE ------------------
	usersRout := router.PathPrefix("/users").Subrouter()
	usersRout.Use(tokenMiddleWare)
	usersRout.Use(formMiddleWare)

	usersRout.HandleFunc("/ping",    ping).Methods("GET")
	usersRout.HandleFunc("/logout",  user_logout).Methods("GET")
	usersRout.HandleFunc("/rename",  user_rename).Methods("PUT")
	usersRout.HandleFunc("/repass",  user_repass).Methods("PUT")
	// -------------------------------------------------

	// ------------------ ROOMS ROUTE ------------------
	roomsRout := router.PathPrefix("/rooms").Subrouter()
	roomsRout.Use(tokenMiddleWare)
	roomsRout.Use(formMiddleWare)

	roomsRout.HandleFunc("/build",   room_build).Methods("POST")
	roomsRout.HandleFunc("/close",   room_close).Methods("DELETE")
	roomsRout.HandleFunc("/rename",  room_rename).Methods("PUT")
	roomsRout.HandleFunc("/publist", pub_list).Methods("GET")
	// -------------------------------------------------

	// #------------------ WEBSOCKET ------------------#
	wsRout := router.PathPrefix("/chat").Subrouter()

	wsRout.HandleFunc("/manage",  rooms_management)
	// #-----------------------------------------------#

	return router
}
