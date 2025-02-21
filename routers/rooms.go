package routers

import (
	"net/http"

	"TermoChat/components"
)


func room_build(w http.ResponseWriter, r *http.Request) {
  name := r.URL.Query().Get("name")
  creator_hash := r.URL.Query().Get("creator_hash")
  is_public_str := r.URL.Query().Get("is_public")

  var room *components.Room
  room.Init(name, creator_hash, is_public_str)

  var db *components.Database
  err := db.BuildRoom(name, creator_hash, room.IsPublic)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
}

func room_close(w http.ResponseWriter, r *http.Request) {
  hash := r.URL.Query().Get("hash")

  var db *components.Database
  db.CloseRoom(hash)
}
