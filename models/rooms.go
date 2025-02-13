package models


type Room struct {
  Name        string
  CreatorHash string
  Hash        string
  IsPublic    bool
  Clients     []*RoomClient
}
