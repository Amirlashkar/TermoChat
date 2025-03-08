package routers

import (
	"encoding/json"
	"net/http"
	"net/url"

    "TermoChat/components"
)


// Checks token validity & saves it on format
func tokenMiddleWare(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        is_valid, _, hash := checkToken(r)

        if !is_valid {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(&map[string]string{"error": "Unauthorized token"})
            return
        }

        // Check if the user is logged in
        var db *components.Database
        user, _ := db.GetUser("", *hash)
        if user.IsLogged == false {
            writeJsonResp(w, http.StatusUnauthorized, "", "Log in first")
            return
        }

        if r.Form == nil {
            r.Form = make(url.Values)
        }

        r.Form.Add("thash", *hash) // make token derived hash accessible

        next.ServeHTTP(w, r)
    })
}

// Check if there is a form & it follows allowed format
func formMiddleWare(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("Content-Type") != "application/json" {
            response := responseBuilder("", "Please enter data as right Content-Type")
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(&response)
            return
        }

        var formData map[string]string
        err := json.NewDecoder(r.Body).Decode(&formData)
        if len(formData) != 0 {
            if err != nil {
                response := responseBuilder("", err.Error())
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(&response)
                return
            }

            if r.Form == nil {
                r.Form = make(url.Values)
            }

            for k, v := range formData {
                r.Form.Add(k, v)
            }
        }
        next.ServeHTTP(w, r)
    })
}
