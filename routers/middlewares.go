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
        var db *components.Database

        is_valid, _, hash := checkToken(r)

        if !is_valid {
            w.Header().Set("Content-Type", "application/json")
            jsonResp(w, "", "Unauthorized token ; login again", "error", http.StatusUnauthorized, map[string]any{})
            return
        }

        user, _ := db.GetUser("", *hash)
        // Check if the user is logged in
        if !user.IsLogged {
            jsonResp(w, "", "Log in first", "error", http.StatusBadRequest, map[string]any{})
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
            jsonResp(w, "", "Please enter data as right Content-Type", "error", http.StatusBadRequest, map[string]any{})
            return
        }

        var formData map[string]string
        err := json.NewDecoder(r.Body).Decode(&formData)
        if len(formData) != 0 {
            if err != nil {
                jsonResp(w, "", err.Error(), "error", http.StatusBadRequest, map[string]any{})
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
