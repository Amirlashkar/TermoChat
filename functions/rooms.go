package functions


import (
	"termo_chat/models"
)


// Adds client to some room
func AddClient(client *models.RoomClient, room *models.Room) *models.Room {
  room.Clients = append(room.Clients, client)
  return room
}

// Removes a client from a sample room
func RemoveClient(client *models.RoomClient, room *models.Room) *models.Room {
  var newClients []*models.RoomClient
  for _, ptr := range room.Clients {
    if ptr != client {
      newClients = append(newClients, ptr)
    }
  }

  newRoom := models.Room {
    Hash: room.Hash,
    Clients: newClients,
  }
  return &newRoom
}


