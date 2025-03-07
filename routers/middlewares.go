package routers


import (
    "net/http"
    "encoding/json"
)


func authMiddleWare(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        is_valid, _ := checkToken(r)
        if !is_valid {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized token"})
            return
        }
        next.ServeHTTP(w, r)
    })
}

func formMiddleWare(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        err := r.ParseForm()
        if err != nil {
            response := responseBuilder("", err)
            json.NewEncoder(w).Encode(&response)
            return
        }
        next.ServeHTTP(w, r)
    })
}
