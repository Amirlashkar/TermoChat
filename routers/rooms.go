package routers

import (
	"fmt"
	"net/http"

	"TermoChat/components"
)


func room_build(w http.ResponseWriter, r *http.Request) {
    name             := r.FormValue("name")
    creator_hash     := r.FormValue("thash")
    is_public_str    := r.FormValue("is_public_str")

    var room *components.Room
    room.Init(name, creator_hash, is_public_str)

    var db *components.Database
    err := db.BuildRoom(name, creator_hash, room.IsPublic)
    if err != nil {
		jsonResp(w, "", err.Error(), "error", http.StatusBadRequest, map[string]any{})
        return
    }
    message := fmt.Sprintf("Room %s built", room.Name)
    jsonResp(w, message, "", "ok", http.StatusOK, map[string]any{})
}

func room_close(w http.ResponseWriter, r *http.Request) {
    user_hash     := r.FormValue("thash")
    hash          := r.FormValue("hash")

    var db *components.Database
    err := db.CloseRoom(user_hash, hash)
    if err != nil {
		jsonResp(w, "", err.Error(), "error", http.StatusBadRequest, map[string]any{})
        return
    }
    message := fmt.Sprintf("Room %s deleted", hash)
    jsonResp(w, message, "", "ok", http.StatusOK, map[string]any{})
}

func room_rename(w http.ResponseWriter, r *http.Request) {
    user_hash     := r.FormValue("thash")
    hash          := r.FormValue("hash")
    new_name      := r.FormValue("new_name")

    var db *components.Database
    room, err := db.GetRoom("", hash)
    if err != nil {
        jsonResp(w, "", err.Error(), "error", http.StatusBadRequest, map[string]any{})
        return
    }

    room.Name = new_name
    room.ReInit() // to update room hash
    err = db.UpdateRoom(user_hash, room, hash)
    if err != nil {
        jsonResp(w, "", err.Error(), "error", http.StatusBadRequest, map[string]any{})
        return
    }
    message := fmt.Sprintf("Room %s renamed", room.Name)
    jsonResp(w, message, "", "ok", http.StatusOK, map[string]any{})
}

func pub_list(w http.ResponseWriter, r *http.Request) {
    var db *components.Database
    names, hashes, err:= db.ListRooms()
    if err != nil {
        jsonResp(w, "", err.Error(), "error", http.StatusBadRequest, map[string]any{})
        return
    }

    data := map[string]any {
        "names":  names,
        "hashes": hashes,
    }
    jsonResp(w, "", "", "ok", http.StatusOK, data)
}
