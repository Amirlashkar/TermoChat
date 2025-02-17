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


// Initialize
func (u *User) NewUser(show_name string, password string, related_question string, related_answer string) (*User, error) {
  pass_hash := universal.CreateHash(map[string]interface{} {
    "Password": password,
  })
  user_hash := universal.CreateHash(map[string]interface{}{
    "ShowName":         show_name,
    "PassHash":         pass_hash,
    "RelatedQuestion":  related_question,
    "RelatedAnswer":    related_answer,
  })

  user := &User {
    ShowName:           show_name,
    PassHash:           pass_hash,
    RelatedQuestion:    related_question,
    RelatedAnswer:      related_answer,
    Hash:               user_hash,
    IsLogged:           false,
  }

  _, err := db.GetUser(user.ShowName)
  if err != nil { // if user with such name doesn't exist
    return user, nil
  } else {
    return nil, err
  }
}

// Rename user
func (u *User) Rename(new_name string) {
  u.ShowName = new_name
  u.Hash = universal.CreateHash(map[string]interface{}{
    "ShowName":         new_name,
    "PassHash":         u.PassHash,
    "RelatedQuestion":  u.RelatedQuestion,
    "RelatedAnswer":    u.RelatedAnswer,
  })
  db.UpdateUser(u, u.ShowName)
}

// Changing password
func (u *User) Repass(current_pass string, new_pass string) error {
  current_hash := universal.CreateHash(map[string]interface{}{
    "Password": current_pass,
  })

  if current_hash != u.PassHash {
    return fmt.Errorf("ERROR: Current provided password is wrong")
  } else {
    u.PassHash = universal.CreateHash(map[string]interface{} {
      "Password": new_pass,
    })

    u.Hash = universal.CreateHash(map[string]interface{}{
      "ShowName":         u.ShowName,
      "PassHash":         u.PassHash,
      "RelatedQuestion":  u.RelatedQuestion,
      "RelatedAnswer":    u.RelatedAnswer,
    })
    db.UpdateUser(u, u.ShowName)
    return nil
  }
}

// Making RoomClient type from User type
func (u *User) User2Client(user User, conn *websocket.Conn, room Room) *RoomClient {
  var client = RoomClient {
    UserHash: user.Hash,
    Conn:     conn,
    RoomHash: room.Hash,
  }
  return &client
}
