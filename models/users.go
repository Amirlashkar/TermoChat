package models


import (
	"github.com/gorilla/websocket"
)


type User struct {
  ShowName        string
  PassHash        string
  RelatedQuestion string
  RelatedAnswer   string
  Hash            string
  IsLogged        bool
}

type RoomClient struct {
  UserHash  string
  Conn      *websocket.Conn
  RoomHash  string
}
