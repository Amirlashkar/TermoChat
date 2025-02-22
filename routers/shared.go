package routers


import (
  "sync"
  "TermoChat/components"
)


var mu sync.Mutex
var rooms map[string]map[string]*components.RoomClient
