package routers

import (
	"fmt"

	"TermoChat/universal"

	"github.com/gorilla/mux"
)

// Main function to provide last router
func ProvideRouter() *mux.Router{
  router := mux.NewRouter()

  // ------------------ USERS ROUT ------------------
  usersH := universal.Word2Hash("users")
  usersURL := fmt.Sprintf("/%s", usersH)
  usersRout := router.PathPrefix(usersURL).Subrouter()

  signupH := universal.Word2Hash("signup")
  signupURL := fmt.Sprintf("/%s", signupH)
  loginH := universal.Word2Hash("login")
  loginURL := fmt.Sprintf("/%s", loginH)
  logoutH := universal.Word2Hash("logout")
  logoutURL := fmt.Sprintf("/%s", logoutH)
  renameH := universal.Word2Hash("rename")
  renameURL := fmt.Sprintf("/%s", renameH)
  repassH := universal.Word2Hash("repass")
  repassURL := fmt.Sprintf("/%s", repassH)

  usersRout.HandleFunc(signupURL, user_signup).Methods("GET")
  usersRout.HandleFunc(loginURL, user_login).Methods("GET")
  usersRout.HandleFunc(logoutURL, user_logout).Methods("GET")
  usersRout.HandleFunc(renameURL, user_rename).Methods("GET")
  usersRout.HandleFunc(repassURL, user_repass).Methods("GET")
  // ------------------------------------------------

  // ------------------ ROOMS ROUT ------------------
  roomsH := universal.Word2Hash("rooms")
  roomsURL := fmt.Sprintf("/%s", roomsH)
  roomsRout := router.PathPrefix(roomsURL).Subrouter()

  // ------------------ WEBSOCKET ------------------
  manageH := universal.Word2Hash("manage")
  manageURL := fmt.Sprintf("/%s", manageH)

  roomsRout.HandleFunc(manageURL, rooms_management)
  // ------------------------------------------------
  buildH := universal.Word2Hash("build")
  buildURL := fmt.Sprintf("/%s", buildH)
  closeH := universal.Word2Hash("close")
  closeURL := fmt.Sprintf("/%s", closeH)

  roomsRout.HandleFunc(buildURL, room_build).Methods("GET")
  roomsRout.HandleFunc(closeURL, room_close).Methods("GET")
  roomsRout.HandleFunc(renameURL, room_rename).Methods("GET")
  // ------------------------------------------------

  return router
}
