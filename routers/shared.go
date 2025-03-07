package routers

import (
	"net/http"
	"strings"
	"sync"

	"TermoChat/components"
	"TermoChat/universal"
)


var mu sync.Mutex
var rooms map[string]map[string]*components.RoomClient


func CheckToken(r *http.Request) (bool, error) {
    token := r.Header.Get("Authorization")
    token = strings.TrimPrefix(token, "Bearer ")

    return universal.IsTokenValid(token)
}
