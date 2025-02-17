package components


import (
	"sync"
  "TermoChat/universal"
)


type Room struct {
  Name        string
  CreatorHash string
  Hash        string
  IsPublic    bool
  Clients     []string
}


// Initialize room
func (r *Room) NewRoom(name string, creator_hash string, is_public bool) (*Room, error) {
  hash := universal.CreateHash(map[string]interface{} {
    "Name":       name,
    "CreateHash": creator_hash,
  })

  room := &Room {
    Name:         name,
    CreatorHash:  creator_hash,
    Hash:         hash,
    IsPublic:     is_public,
    Clients:      []string{creator_hash},
  }

  _, err := db.GetRoom(room.Name)
  if err != nil {
    return room, nil
  } else {
    return nil, err
  }
}

// Adds client to a room
func (r *Room) AddClient(client RoomClient) {
  mu.Lock()
  defer mu.Unlock()

  r.Clients = append(r.Clients, client.UserHash)
  db.UpdateRoom(r, r.Name)
}

// Removes a client from a room
func (r *Room) RemoveClient(client RoomClient) {
  mu.Lock()
  defer mu.Unlock()

  for i, ptr := range r.Clients {
    if ptr == client.UserHash {
      r.Clients = append(r.Clients[:i], r.Clients[i+1:]...)
    }
  }
  db.UpdateRoom(r, r.Name)
}
