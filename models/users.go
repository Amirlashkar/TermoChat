package models


import (
	"github.com/gorilla/websocket"
)


type User struct {
  Showname string
  Password string
  RelatedQuestion string
  RelatedAnswer string
  Hash string
  is_logged bool
}

type RoomClient struct {
  User *User
  Conn *websocket.Conn
  RoomHash string
}
