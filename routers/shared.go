package routers

import (
	"sync"
	"strings"
	"net/http"
    "encoding/json"

	"TermoChat/components"
	"TermoChat/universal"
)


var mu sync.Mutex
var rooms map[string]map[string]*components.RoomClient
type response struct {
    message   string `json:"message"`
    error     error  `json:"error"`
}

func checkToken(r *http.Request) (bool, error) {
    token := r.Header.Get("Authorization")
    token = strings.TrimPrefix(token, "Bearer ")

    return universal.IsTokenValid(token)
}

func responseBuilder(message string, err error) response {
    return response {
        message: message,
        error:   err,
    }
}

func writeJsonResp(w http.ResponseWriter, status_code int, message string, err error) {
        w.WriteHeader(status_code)
        json.NewEncoder(w).Encode(
            responseBuilder(message, err),
        )
}
