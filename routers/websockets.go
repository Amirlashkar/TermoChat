package routers


import (
  "log"
  "net/http"

  "TermoChat/components"
	"github.com/gorilla/websocket"
)


var rooms = map[string]map[string]*components.RoomClient {}

var upgrader = websocket.Upgrader {
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections
	},
}

// Broadcast a supposed message to all clients on same room
func broadcast(room map[string]*components.RoomClient, message []byte, db *components.Database, room_hash string) {
  for client := range room {
    err := room[client].Conn.WriteMessage(websocket.TextMessage, message)
    if err != nil {
      log.Printf("SUCCESS: User(%s) left room(%s)", client, room_hash)
      room[client].Conn.Close()
      room_struct, _ := db.GetRoom("", room_hash)

      room_struct.RemoveClient(*room[client])
      db.UpdateRoom(room_struct, room_hash)
      delete(room, client)
    }
  }
}

// Manages each websocket connection that a client initializes
func rooms_management(w http.ResponseWriter, r *http.Request) {
  user_hash := r.URL.Query().Get("user_hash")
  room_hash := r.URL.Query().Get("room_hash")

  var db components.Database
  user, err := db.GetUser("", user_hash)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  room, err := db.GetRoom("", room_hash)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  conn, _ := upgrader.Upgrade(w, r, nil)
  defer conn.Close()

  client := components.RoomClient{
    UserHash: user_hash,
    Conn:     conn,
    RoomHash: room_hash,
  }

  room.AddClient(client)
  db.UpdateRoom(room, room.Hash)
  log.Printf("SUCCESS: User(%s) joined room(%s)", user.Hash, room.Hash)

  rooms[room.Hash][client.UserHash] = &client

  for {
    _, message, err := conn.ReadMessage()
    if err != nil {
      log.Printf("SUCCESS: User(%s) left room(%s)", user.Hash, room.Hash)
      room.RemoveClient(client)
      db.UpdateRoom(room, room.Hash)
    } else {
      broadcast(rooms[room.Hash], message, &db, room.Hash)
    }
  }
}
