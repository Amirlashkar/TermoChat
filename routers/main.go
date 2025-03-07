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

	authRout.HandleFunc("/signup", user_signup).Methods("POST")
	authRout.HandleFunc("/login",  user_login).Methods("POST")
	// ------------------------------------------


	// ------------------ USERS ROUTE ------------------
	usersRout := router.PathPrefix("/users").Subrouter()
    usersRout.Use(authMiddleWare)
    usersRout.Use(formMiddleWare)

	usersRout.HandleFunc("/logout", user_logout).Methods("POST")
	usersRout.HandleFunc("/rename", user_rename).Methods("PUT")
	usersRout.HandleFunc("/repass", user_repass).Methods("PUT")
	// -------------------------------------------------


	// ------------------ ROOMS ROUTE ------------------
	roomsRout := router.PathPrefix("/rooms").Subrouter()
    roomsRout.Use(authMiddleWare)
    roomsRout.Use(formMiddleWare)

	// #------------------ WEBSOCKET ------------------#
	roomsRout.HandleFunc("/manage", rooms_management)
	// #-----------------------------------------------#

	roomsRout.HandleFunc("/build",  room_build).Methods("POST")
	roomsRout.HandleFunc("/close",  room_close).Methods("DELETE")
    roomsRout.HandleFunc("/rename", room_rename).Methods("PUT")
    // roomsRout.HandleFunc("/list_rooms", list_rooms).Methods("GET")
	// -------------------------------------------------


	return router
}
