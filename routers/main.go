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

  usersRout := router.PathPrefix(usersURL).Subrouter()
  usersRout.HandleFunc(signupURL, signup).Methods("GET")
  usersRout.HandleFunc(loginURL, login).Methods("GET")
  usersRout.HandleFunc(logoutURL, logout).Methods("GET")
  usersRout.HandleFunc(renameURL, rename).Methods("GET")
  usersRout.HandleFunc(repassURL, repass).Methods("GET")
  // ------------------------------------------------

  // ------------------ ROOMS ROUT ------------------
  // ------------------------------------------------

  return router
}
