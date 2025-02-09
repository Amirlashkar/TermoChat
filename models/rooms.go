package models


type Room struct {
  Hash string
  Clients []*RoomClient
}
