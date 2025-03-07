package routers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"TermoChat/components"
	"TermoChat/universal"
)

func user_signup(w http.ResponseWriter, r *http.Request) {
    showname         := r.FormValue("showname")
    password         := r.FormValue("password")
    related_question := r.FormValue("related_question")
    related_answer   := r.FormValue("related_answer")

	var user *components.User
	user.Init(showname, password, related_question, related_answer)

	var db *components.Database
    err := db.SignUp(user)
	if err != nil {
        writeJsonResp(w, http.StatusBadRequest, "", err)
	} else {
        message := fmt.Sprintf("User %s signed up successfully", user.Hash)
        writeJsonResp(w, http.StatusAccepted, message, nil)
    }
}

func user_login(w http.ResponseWriter, r *http.Request) {
    showname         := r.FormValue("showname")
    password         := r.FormValue("password")

	var db *components.Database
	user, err := db.LogIn(showname, password)
	if err != nil {
        writeJsonResp(w, http.StatusBadRequest, "", err)
	} else {
		token_resp := universal.GenerateJWT(user.Hash, 1*time.Hour)
        w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(&token_resp)
	}
}

func user_logout(w http.ResponseWriter, r *http.Request) {
    hash := r.FormValue("hash")

	var db *components.Database
    err := db.LogOut(hash)
	if err != nil {
        writeJsonResp(w, http.StatusBadRequest, "", err)
	} else {
        message := fmt.Sprintf("User %s logged out", hash)
        writeJsonResp(w, http.StatusOK, message, nil)
    }
}

func user_rename(w http.ResponseWriter, r *http.Request) {
    hash     := r.FormValue("hash")
    new_name := r.FormValue("new_name")

    var db *components.Database
    user, err := db.GetUser("", hash)
    if err != nil {
        writeJsonResp(w, http.StatusBadRequest, "", err)
    } else {
        user.ShowName = new_name
        user.ReInit()
        db.UpdateUser(user, hash)
        message := fmt.Sprintf("User %s renamed ; NewHash: %s", hash, user.Hash)
        writeJsonResp(w, http.StatusOK, message, nil)
    }
}

func user_repass(w http.ResponseWriter, r *http.Request) {
    hash         := r.FormValue("hash")
    current_pass := r.FormValue("current_pass")
    new_pass     := r.FormValue("new_pass")

	var db *components.Database
	user, err := db.GetUser("", hash)
	if err != nil {
        writeJsonResp(w, http.StatusBadRequest, "", err)
	} else {
		err := user.Repass(current_pass, new_pass)
		if err != nil {
            writeJsonResp(w, http.StatusUnauthorized, "", err)
		} else {
			db.UpdateUser(user, hash)
            message := fmt.Sprintf("User %s password changed", user.Hash)
            writeJsonResp(w, http.StatusOK, message, nil)
		}
	}
}
