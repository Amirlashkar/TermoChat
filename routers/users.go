package routers

import (
	"fmt"
	"net/http"
	"time"

	"TermoChat/components"
	"TermoChat/universal"
)

func user_signup(w http.ResponseWriter, r *http.Request) {
	showname := r.FormValue("show_name")
	password := r.FormValue("password")
	related_question := r.FormValue("related_question")
	related_answer := r.FormValue("related_answer")

	user := &components.User{}
	user.Init(showname, password, related_question, related_answer)

	var db *components.Database
	err := db.SignUp(user)
	if err != nil {
		jsonResp(w, "", err.Error(), "error", http.StatusBadRequest, map[string]any{})
		return
	} else {
		message := fmt.Sprintf("User %s signed up successfully", user.Hash)
		jsonResp(w, message, "", "ok", http.StatusOK, map[string]any{})
		return
	}
}

func user_login(w http.ResponseWriter, r *http.Request) {
	showname := r.FormValue("show_name")
	password := r.FormValue("password")

	var db *components.Database
	user, err := db.LogIn(showname, password)
	if err != nil {
		jsonResp(w, "", err.Error(), "error", http.StatusBadRequest, map[string]any{})
	} else {
		token_resp := universal.GenerateJWT(user.Hash, 1*time.Hour)
        data := map[string]any {
            "token":      token_resp.Token,
            "expires_at": token_resp.ExpiresAt,
        }
		jsonResp(w, "User logged in", "", "ok", http.StatusOK, data)
	}
}

func user_existance(w http.ResponseWriter, r *http.Request) {
    show_name := r.FormValue("show_name")

	var db *components.Database
    user, err := db.GetUser("show_name", show_name)
    if user == nil {
		jsonResp(w, "", err.Error(), "error", http.StatusBadRequest, map[string]any{})
    } else {
		jsonResp(w, "user exists", "", "ok", http.StatusOK, map[string]any{})
    }
}

// We use this function after tokenMiddleWare
// to check if token is valid or user connection
// would cause any issue
func ping(w http.ResponseWriter, r *http.Request) {
    jsonResp(w, "PONG", "", "ok", http.StatusOK, map[string]any{})
}

func user_logout(w http.ResponseWriter, r *http.Request) {
	thash := r.FormValue("thash")

	var db *components.Database
	err := db.LogOut(thash)
	if err != nil {
        jsonResp(w, "", err.Error(), "error", http.StatusBadRequest, map[string]any{})
	} else {
		message := fmt.Sprintf("User %s logged out", thash)
        jsonResp(w, message, "", "ok", http.StatusOK, map[string]any{})
	}
}

func user_rename(w http.ResponseWriter, r *http.Request) {
	hash := r.FormValue("thash")
	new_name := r.FormValue("new_name")

	var db *components.Database
	user, err := db.GetUser("", hash)
	if err != nil {
        jsonResp(w, "", err.Error(), "error", http.StatusBadRequest, map[string]any{})
	} else {
		user.ShowName = new_name
		db.UpdateUser(user, hash)
		message := fmt.Sprintf("User %s renamed ; NewName: %s", hash, user.ShowName)
        jsonResp(w, message, "", "ok", http.StatusOK, map[string]any{})
	}
}

func user_repass(w http.ResponseWriter, r *http.Request) {
	hash := r.FormValue("thash")
	current_pass := r.FormValue("current_pass")
	new_pass := r.FormValue("new_pass")

	var db *components.Database
	user, err := db.GetUser("", hash)
	if err != nil {
        jsonResp(w, "", err.Error(), "error", http.StatusBadRequest, map[string]any{})
	} else {
		err := user.Repass(current_pass, new_pass)
		if err != nil {
            jsonResp(w, "", err.Error(), "error", http.StatusUnauthorized, map[string]any{})
		} else {
			db.UpdateUser(user, hash)
			message := fmt.Sprintf("User %s password changed", user.ShowName)
            jsonResp(w, message, "", "error", http.StatusOK, map[string]any{})
		}
	}
}
