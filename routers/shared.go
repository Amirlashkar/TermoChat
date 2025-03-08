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

type RoomsList struct {
    Names     []string `json:"names"`
    Hashes    []string `json:"hashes"`
}

type response struct {
    Message   string   `json:"message"`
    Error     string   `json:"error,omitempty"`
}

func checkToken(r *http.Request) (bool, error, *string) {
    token := r.Header.Get("Authorization")
    token = strings.TrimPrefix(token, "Bearer ")

    return universal.IsTokenValid(token)
}

func responseBuilder(message string, err string) *response {
    return &response {
        Message: message,
        Error:   err,
    }
}

func writeJsonResp(w http.ResponseWriter, status_code int, message string, err string) {
        w.WriteHeader(status_code)
        json.NewEncoder(w).Encode(
            responseBuilder(message, err),
        )
}
