package routers

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"TermoChat/components"
	"TermoChat/universal"
)


var mu sync.Mutex
var rooms map[string]map[string]*components.RoomClient

type RoomsList struct {
    Names     []string `json:"names"`
    Hashes    []string `json:"hashes"`
}

func checkToken(r *http.Request) (bool, error, *string) {
    token := r.Header.Get("Authorization")
    token = strings.TrimPrefix(token, "Bearer ")

    return universal.IsTokenValid(token)
}

type response struct {
    Message   string           `json:"message,omitempty"`
    Error     string           `json:"error,omitempty"`
    Status    string           `json:"status"`
    Code      int              `json:"code"`
    Data      map[string]any   `json:"data,omitempty"`
}

func responseBuilder(message string, err string, status string,
                     code int, data map[string]any) *response {
    return &response {
        Message: message,
        Error:   err,
        Status:  status,
        Code:    code,
        Data:    data,
    }
}

func jsonResp(w http.ResponseWriter, message string, err string, status string, code int, data map[string]any) {
        w.WriteHeader(code)
        json.NewEncoder(w).Encode(
            responseBuilder(message, err, status, code, data),
        )
}

func addClient2Room(room_hash string, client *components.RoomClient) {
    if rooms == nil {
        rooms = make(map[string]map[string]*components.RoomClient)
    }

    if rooms[room_hash] == nil {
        rooms[room_hash] = make(map[string]*components.RoomClient)
    }

    rooms[room_hash][client.UserHash] = client
}
