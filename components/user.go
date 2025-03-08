package components


import (
  "fmt"

	"TermoChat/universal"

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
  UserHash        string
  Conn            *websocket.Conn
  RoomHash        string
}


// Initialize new user
func (u *User) Init(show_name string, password string, related_question string, related_answer string) {
    u.PassHash = universal.Data2Hash(map[string]any {
        "Password": password,
    })
    u.Hash = universal.Data2Hash(map[string]any{
        "ShowName":         show_name,
        "PassHash":         u.PassHash,
        "RelatedQuestion":  related_question,
        "RelatedAnswer":    related_answer,
    })

    u.ShowName         = show_name
    u.RelatedQuestion  = related_question
    u.RelatedAnswer    = related_answer
    u.IsLogged         = false
}

// Reinitializes to update user hash due to changed credentials
func (u *User) ReInit() {
  u.Hash = universal.Data2Hash(map[string]any {
    "ShowName":         u.ShowName,
    "PassHash":         u.PassHash,
    "RelatedQuestion":  u.RelatedQuestion,
    "RelatedAnswer":    u.RelatedAnswer,
  })
}

// Changing password
func (u *User) Repass(current_pass string, new_pass string) error {
  current_hash := universal.Data2Hash(map[string]any {
    "Password": current_pass,
  })

  if current_hash != u.PassHash {
    return fmt.Errorf("Current provided password is wrong")
  } else {
    u.PassHash = universal.Data2Hash(map[string]any {
      "Password": new_pass,
    })

    // u.ReInit() // If this happens then we're up to change token per put requests
    return nil
  }
}

// Making RoomClient type from User type
func (u *User) User2Client(conn *websocket.Conn, room *Room) *RoomClient {
  var client = RoomClient {
    UserHash: u.Hash,
    Conn:     conn,
    RoomHash: room.Hash,
  }
  return &client
}
