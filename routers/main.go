package routers

import (
	"github.com/gorilla/mux"
)

// Main function to provide last router
func ProvideRouter() *mux.Router {
	router := mux.NewRouter()

	// ------------------ USERS ROUTE ------------------
	usersRout := router.PathPrefix("/users").Subrouter()

	usersRout.HandleFunc("/signup", user_signup).Methods("GET")
	usersRout.HandleFunc("/login", user_login).Methods("GET")
	usersRout.HandleFunc("/logout", user_logout).Methods("GET")
	usersRout.HandleFunc("/rename", user_rename).Methods("GET")
	usersRout.HandleFunc("/repass", user_repass).Methods("GET")
	// -------------------------------------------------

	// ------------------ ROOMS ROUTE ------------------
	roomsRout := router.PathPrefix("/rooms").Subrouter()

	// ------------------ WEBSOCKET ------------------
	roomsRout.HandleFunc("/manage", rooms_management)
	// ------------------------------------------------

	roomsRout.HandleFunc("/build", room_build).Methods("GET")
	roomsRout.HandleFunc("/close", room_close).Methods("GET")
    roomsRout.HandleFunc("/rename", room_rename).Methods("GET")
	// -------------------------------------------------

	return router
}
