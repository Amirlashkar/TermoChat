package components


import (
  "TermoChat/universal"
)


type Room struct {
  Name        string
  CreatorHash string
  Hash        string
  IsPublic    bool
  Clients     []string
}


// Initialize new room
func (r *Room) Init(name string, creator_hash string, is_public bool) {
  r.Hash = universal.Data2Hash(map[string]interface{} {
    "Name":         name,
    "CreatorHash":  creator_hash,
  })

  r = &Room {
    Name:         name,
    CreatorHash:  creator_hash,
    Hash:         r.Hash,
    IsPublic:     is_public,
    Clients:      []string{creator_hash},
  }
}

// Update room hash due to updated room details
func (r *Room) ReInit() {
  r.Hash = universal.Data2Hash(map[string]interface{} {
    "Name":         r.Name,
    "CreatorHash":  r.CreatorHash,
  })
}

// Adds client to a room
func (r *Room) AddClient(client RoomClient) {
  mu.Lock()
  defer mu.Unlock()

  r.Clients = append(r.Clients, client.UserHash)
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
}
