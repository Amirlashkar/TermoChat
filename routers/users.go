package routers


import (
	"time"
	"encoding/json"
	"net/http"

	"TermoChat/components"
	"TermoChat/universal"
)


func user_signup(w http.ResponseWriter, r *http.Request) {
  showname := r.URL.Query().Get("showname")
  password := r.URL.Query().Get("password")
  related_question := r.URL.Query().Get("related_question")
  related_answer := r.URL.Query().Get("related_answer")

  var user *components.User
  user.Init(showname, password, related_question, related_answer)

  var db *components.Database
  err := db.SignUp(user)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
  }
}

func user_login(w http.ResponseWriter, r *http.Request) {
  // Taking credentials from queries
  showname := r.URL.Query().Get("showname")
  password := r.URL.Query().Get("password")

  var db *components.Database
  user, err := db.LogIn(showname, password)
  if err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
  } else {
    token := universal.GenerateJWT(user.Hash, 1*time.Hour)
    json.NewEncoder(w).Encode(&token)
  }
}

func user_logout(w http.ResponseWriter, r *http.Request) {
  hash := r.URL.Query().Get("hash")

  var db *components.Database
  err := db.LogOut(hash)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
  }
}

func user_rename(w http.ResponseWriter, r *http.Request) {
  hash := r.URL.Query().Get("hash")
  new_name := r.URL.Query().Get("new_name")

  var db *components.Database
  user, err := db.GetUser("", hash)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
  } else {
    user.ShowName = new_name
    user.ReInit()
    db.UpdateUser(user, hash)
  }
}

func user_repass(w http.ResponseWriter, r *http.Request) {
  hash := r.URL.Query().Get("hash")
  current_pass := r.URL.Query().Get("current")
  new_pass := r.URL.Query().Get("new")

  var db *components.Database
  user, err := db.GetUser("", hash)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
  } else {
    err := user.Repass(current_pass, new_pass)
    if err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
    } else {
      db.UpdateUser(user, hash)
    }
  }
}
