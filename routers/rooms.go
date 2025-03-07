package routers

import (
	"net/http"

	"TermoChat/components"
)


func room_build(w http.ResponseWriter, r *http.Request) {
    is_valid, token_err := CheckToken(r)
    if !is_valid {
        http.Error(w, token_err.Error(), http.StatusBadRequest)
        return
    }

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
	is_valid, token_err := CheckToken(r)
	if !is_valid {
		http.Error(w, token_err.Error(), http.StatusBadRequest)
		return
	}

    hash := r.URL.Query().Get("hash")

    var db *components.Database
    err := db.CloseRoom(hash)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
    }
}

func room_rename(w http.ResponseWriter, r *http.Request) {
	is_valid, token_err := CheckToken(r)
	if !is_valid {
		http.Error(w, token_err.Error(), http.StatusBadRequest)
		return
	}

    hash := r.URL.Query().Get("hash")
    new_name := r.URL.Query().Get("new_name")

    var db *components.Database
    room, err := db.GetRoom("", hash)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
    }
    room.Name = new_name
    room.ReInit()
    err = db.UpdateRoom(room, hash)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
    }
}
