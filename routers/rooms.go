package routers

import (
	"fmt"
	"net/http"

	"TermoChat/components"
)


func room_build(w http.ResponseWriter, r *http.Request) {
    name             := r.FormValue("name")
    creator_hash     := r.FormValue("creator_hash")
    is_public_str    := r.FormValue("is_public_str")

    var room *components.Room
    room.Init(name, creator_hash, is_public_str)

    var db *components.Database
    err := db.BuildRoom(name, creator_hash, room.IsPublic)
    if err != nil {
        writeJsonResp(w, http.StatusOK, "", err)
        return
    }
    message := fmt.Sprintf("Room %s built", room.Hash)
    writeJsonResp(w, http.StatusCreated, message, nil)
}

func room_close(w http.ResponseWriter, r *http.Request) {
    hash     := r.FormValue("hash")

    var db *components.Database
    err := db.CloseRoom(hash)
    if err != nil {
        writeJsonResp(w, http.StatusOK, "", err)
        return
    }
    message := fmt.Sprintf("Room %s deleted", hash)
    writeJsonResp(w, http.StatusOK, message, nil)
}

func room_rename(w http.ResponseWriter, r *http.Request) {
    hash     := r.FormValue("hash")
    new_name := r.FormValue("new_name")

    var db *components.Database
    room, err := db.GetRoom("", hash)
    if err != nil {
        writeJsonResp(w, http.StatusBadRequest, "", err)
        return
    }
    room.Name = new_name
    room.ReInit() // to update room hash
    err = db.UpdateRoom(room, hash)
    if err != nil {
        writeJsonResp(w, http.StatusBadRequest, "", err)
        return
    }
    message := fmt.Sprintf("Room %s renamed", hash)
    writeJsonResp(w, http.StatusOK, message, nil)
}
