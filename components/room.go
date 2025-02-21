package components


import (
  "strconv"

  "TermoChat/universal"
)


type Room struct {
  Name        string
  CreatorHash string
  Hash        string
  IsPublic    bool
  Clients     []string // User hashes appends here
}


// Initialize new room
func (r *Room) Init(name string, creator_hash string, is_public_str string) {
  r.Hash = universal.Data2Hash(map[string]interface{} {
    "Name":         name,
    "CreatorHash":  creator_hash,
  })

  is_public, err := strconv.ParseBool(is_public_str)
  if err != nil {
    is_public = false
  }

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

  for i, n := range r.Clients {
    if n == client.UserHash {
      r.Clients = append(r.Clients[:i], r.Clients[i+1:]...)
    }
  }
}
