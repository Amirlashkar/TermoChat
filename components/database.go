package components


import (
	"fmt"
	"log"

	"TermoChat/config"
  "TermoChat/universal"

	"database/sql"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)


// Database instance
type Database struct {}

// Simple db connection function
func (DB *Database) DBConnect() *sql.DB{
  env := config.LoadEnv()
  connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", env.DB_USER, env.DB_PASS, env.DB_NAME)
  db, err := sql.Open("postgres", connStr)
  if err != nil {
    log.Fatal(err)
  }
  return db
}

// Database migration
func (DB *Database) Migration() {
  db := DB.DBConnect()
  defer db.Close()

  // Check if wanted tables exist
  execSQL := `SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename = $1`

  row := db.QueryRow(execSQL, "users")
  var table string
  err := row.Scan(&table)
  log.Println(table)
  if err != nil {
    log.Println(err)
  }

  if table != "users" {
    // Make migrations
    execSQL = `
      CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        show_name VARCHAR(255) NOT NULL,
        pass_hash CHAR(64) NOT NULL,
        related_question TEXT NOT NULL,
        related_answer TEXT NOT NULL,
        hash CHAR(64) NOT NULL,
        is_logged BOOLEAN DEFAULT FALSE
      );

      CREATE TABLE IF NOT EXISTS rooms (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        creator_hash CHAR(64) NOT NULL,
        hash CHAR(64) NOT NULL,
        is_public BOOLEAN DEFAULT TRUE,
        clients CHAR(64)[]
      );`

    db.Exec(execSQL)
    log.Println("SUCCESS: New tables migrated successfully")
  }
}

// Searching for a user
func (DB *Database) GetUser(by string, search_word string) (*User, error) {
  db := DB.DBConnect()
  defer db.Close()

  execSQL := `SELECT show_name, pass_hash, related_question, related_answer, hash, is_logged
              FROM users WHERE `

  if by == "show_name" {
    execSQL += "show_name = $1;"
  } else {
    execSQL += "hash = $1;"
  }

  row := db.QueryRow(execSQL, search_word)

  var user User
  err := row.Scan(&user.ShowName, &user.PassHash, &user.RelatedQuestion, &user.RelatedAnswer, &user.Hash, &user.IsLogged)
	if err != nil || user.Hash == "" {
    return nil, fmt.Errorf("No such user")
	}
  return &user, nil
}

// List Users
func (DB *Database) ListUsers() ([]string, []string, error) {
  db := DB.DBConnect()
  defer db.Close()

  execSQL := `
    SELECT show_name, hash FROM users;
  `

  rows, err := db.Query(execSQL)
  defer rows.Close()

  if err != nil {
    return nil, nil, err
  } else {
    var show_names []string
    var hashes []string
    for rows.Next() {
      var show_name string
      var hash string
      if err := rows.Scan(&show_name, &hash); err != nil {
        return nil, nil, err
      }
      show_names = append(show_names, show_name)
      hashes = append(hashes, hash)
    }
    return show_names, hashes, nil
  }
}

// Adding one user to database
func (DB *Database) SignUp(user *User) error {
  db := DB.DBConnect()
  defer db.Close()

  _, err := DB.GetUser("", user.Hash)
  if err != nil {
    execSQL := `INSERT INTO users 
            (show_name, pass_hash, related_question, related_answer, hash)
            VALUES ($1, $2, $3, $4, $5)`

    db.Exec(execSQL, user.ShowName,
            user.PassHash, user.RelatedQuestion,
            user.RelatedAnswer, user.Hash,
    )
    log.Printf("SUCCESS: User(%s) signed up successfully", user.Hash)
    return nil
  } else {
    return fmt.Errorf("Chosen name is already used")
  }
}

// Updating user on database
func (DB *Database) UpdateUser(new_user *User, hash string) error {
  db := DB.DBConnect()
  defer db.Close()

  _, err := DB.GetUser("", hash)
  if err != nil {
    execSQL := `UPDATE users
                SET show_name = $1, pass_hash = $2, related_question = $3, 
                related_answer = $4, hash = $5, is_logged = $6
                WHERE show_name = $7`

    db.Exec(execSQL, new_user.ShowName, 
            new_user.PassHash, new_user.RelatedQuestion, 
            new_user.RelatedAnswer, new_user.Hash, new_user.IsLogged, new_user.ShowName)

    return nil
  } else {
    return err
  }
}

// User activation
func (DB *Database) LogIn(show_name string, password string) (*User, error){
  PassHash := universal.Data2Hash(map[string]interface{} {
    "Password": password,
  })

  user, err := DB.GetUser("show_name", show_name)
  if err != nil {
    return nil, err
  } else {
    if PassHash != user.PassHash {
      return nil, fmt.Errorf("Wrong password")
    } else {
      if user.IsLogged {
        return user, fmt.Errorf("User already logged in")
      } else {
        user.IsLogged = true
        DB.UpdateUser(user, user.ShowName)
        log.Printf("SUCCESS: User(%s) logged in", user.Hash)
        return user, nil
      }
    }
  }
}

// User deactivation
func (DB *Database) LogOut(hash string) error {
  user, err := DB.GetUser("", hash)
  if err != nil {
    return err
  } else if user.IsLogged != false {
    return fmt.Errorf("Log in first")
  } else {
    user.IsLogged = false
    DB.UpdateUser(user, user.ShowName)
    log.Printf("SUCCESS: User(%s) logged out", hash)
    return nil
  }
}

// Delete a user data
func (DB *Database) DelUser(hash string) error {
  user, err := DB.GetUser("", hash)
  if err != nil {
    return err
  } else {
    db := DB.DBConnect()
    defer db.Close()

    execSQL := `DELETE FROM users WHERE showname = $1`

    db.Exec(execSQL, user.ShowName)
    log.Printf("SUCCESS: User(%s) deleted", hash)
    return nil
  }
}

// List rooms
func (DB *Database) ListRooms() ([]string, []string, error) {
  db := DB.DBConnect()
  defer db.Close()

  execSQL := `
    SELECT name, hash FROM rooms
  `

  rows, err := db.Query(execSQL)
  defer rows.Close()

  if err != nil {
    return nil, nil, err
  } else {
    var names []string
    var hashes []string
    for rows.Next() {
      var name string
      var hash string
      if err := rows.Scan(&name, &hash); err != nil {
        log.Println("ERROR: ", err)
        return nil, nil, err
      }
      names = append(names, name)
      hashes = append(hashes, hash)
    }
    return names, hashes, nil
  }
}

// Getting room info
func (DB *Database) GetRoom(by string, search_word string) (*Room, error) {
  db := DB.DBConnect()
  defer db.Close()

  execSQL := `
    SELECT name, creator_hash, hash, is_public, clients FROM rooms WHERE 
  `

  if by == "name" {
    execSQL += "name = $1;"
  } else {
    execSQL += "hash = $1;"
  }

  row := db.QueryRow(execSQL, search_word)
  var room Room
  err := row.Scan(
    &room.Name, &room.CreatorHash, &room.Hash,
    &room.IsPublic, pq.Array(&room.Clients),
  )
  if err != nil || room.Name == "" {
    return nil, err
  }

  return &room, nil
}

// Creates new room
func (DB *Database) BuildRoom(name string, creator_hash string, is_public bool) error {
  db := DB.DBConnect()
  defer db.Close()

  roomHash := universal.Data2Hash(map[string]interface{} {
    "Name":         name,
    "CreatorHash":  creator_hash,
  })

  room, err := DB.GetRoom("name", name)
  if err != nil {
    room = &Room {
      Name:         name,
      CreatorHash:  creator_hash,
      Hash:         roomHash,
      IsPublic:     is_public,
      Clients:      []string{creator_hash},
    }

    execSQL := `
      INSERT INTO rooms (
        name, creator_hash, hash, is_public, clients
      ) VALUES (
        $1, $2, $3, $4, $5
      );
    `

    db.Exec(execSQL, &room.Name, &room.CreatorHash,
            &room.Hash, &room.IsPublic, pq.Array(&room.Clients),
    )
    log.Printf("SUCCESS: Room(%s) is built", roomHash)
    return nil
  } else {
    return err
  }
}

// Switches room privacy
func (DB *Database) SwitchRoomPriv(hash string) {
  db := DB.DBConnect()
  defer db.Close()

  room, _ := DB.GetRoom("", hash)
  if room == nil {
    log.Printf("ERROR: Room(%s) doesn't exist", hash)
    return
  }

  is_public := room.IsPublic
  if is_public {
    execSQL := `UPDATE rooms SET is_public = false WHERE hash = $1`
    db.Exec(execSQL, hash)
  } else {
    execSQL := `UPDATE rooms SET is_public = true WHERE hash = $1`
    db.Exec(execSQL, hash)
  }
  log.Printf("SUCCESS: Room(%s) privacy status changed", hash)
}

// Close room
func (DB *Database) CloseRoom(hash string) error {
  db := DB.DBConnect()
  defer db.Close()

  _, err := DB.GetRoom("", hash)
  if err != nil {
    log.Printf("ERROR: Room(%s) doesn't exist", hash)
    return err
  }

  execSQL := `DELETE FROM rooms WHERE hash = $1`
  db.Exec(execSQL, hash)
  log.Printf("SUCCESS: Room(%s) closed", hash)
  return nil
}

// Updating room
func (DB *Database) UpdateRoom(new_room *Room, hash string) error {
  db := DB.DBConnect()
  defer db.Close()

  _, err := DB.GetRoom("", hash)
  if err != nil {
    return err
  } else {
    execSQL := `UPDATE rooms
                SET name = $1, creator_hash = $2, hash = $3, 
                is_public = $4, clients = $5
                WHERE hash = $6`

    db.Exec(execSQL, new_room.Name, 
            new_room.CreatorHash, new_room.Hash, 
            new_room.IsPublic, pq.Array(new_room.Clients),
            hash,)

    log.Printf("SUCCESS: Room(%s) updated", hash)
    return nil
  }
}
