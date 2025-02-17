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
      CREATE TABLE users (
        id SERIAL PRIMARY KEY,
        show_name VARCHAR(255) NOT NULL,
        pass_hash CHAR(64) NOT NULL,
        related_question TEXT NOT NULL,
        related_answer TEXT NOT NULL,
        hash CHAR(64) NOT NULL,
        is_logged BOOLEAN DEFAULT FALSE
      );

      CREATE TABLE rooms (
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
func (DB *Database) GetUser(show_name string) (*User, error) {
  db := DB.DBConnect()
  defer db.Close()

  execSQL := `SELECT show_name, pass_hash, related_question, related_answer, hash, is_logged
              FROM users WHERE show_name = $1`
  row := db.QueryRow(execSQL, show_name)

  var user User
  err := row.Scan(&user.ShowName, &user.PassHash, &user.RelatedQuestion, &user.RelatedAnswer, &user.Hash, &user.IsLogged)
	if err != nil {
		if err == sql.ErrNoRows {
      return nil, fmt.Errorf("ERROR: User(%s) not found", show_name)
		}
		return nil, err
	}
  return &user, nil
}

// List Users
func (DB *Database) ListUsers() ([]string, []string) {
  db := DB.DBConnect()
  defer db.Close()

  execSQL := `
    SELECT show_name, hash FROM users;
  `

  rows, err := db.Query(execSQL)
  defer rows.Close()

  var show_names []string
  var hashes []string
  for rows.Next() {
    var show_name string
    var hash string
    if err := rows.Scan(&show_name, &hash); err != nil {
      log.Println("ERROR: ", err)
      return nil, nil
    }
    show_names = append(show_names, show_name)
    hashes = append(hashes, hash)
  }

  if err != nil {
    log.Printf("ERROR: %s", err)
    return nil, nil
  }
  return show_names, hashes
}

// Adding one user to database
func (DB *Database) SignUp(user *User) {
  db := DB.DBConnect()
  defer db.Close()

  execSQL := `INSERT INTO users 
              (show_name, pass_hash, related_question, related_answer, hash)
              VALUES ($1, $2, $3, $4, $5)`

  _, err := db.Exec(execSQL, user.ShowName,
                    user.PassHash, user.RelatedQuestion,
                    user.RelatedAnswer, user.Hash,
  )

  if err == nil {
    log.Printf("SUCCESS: User(%s) signed up", user.ShowName)
  }
}

// Updating user on database
func (DB *Database) UpdateUser(new_user *User, show_name string) {
  db := DB.DBConnect()
  defer db.Close()

	execSQL := `UPDATE users
              SET show_name = $1, pass_hash = $2, related_question = $3, 
              related_answer = $4, hash = $5, is_logged = $6
              WHERE show_name = $7`

  db.Exec(execSQL, new_user.ShowName, 
          new_user.PassHash, new_user.RelatedQuestion, 
          new_user.RelatedAnswer, new_user.Hash, new_user.IsLogged, new_user.ShowName)

  log.Printf("SUCCESS: User(%s) updated", new_user.ShowName)
}

// Makes user active
func (DB *Database) LogIn(show_name string, password string) (*User, error){
  PassHash := universal.CreateHash(map[string]interface{} {
    "Password": password,
  })

  user, err := DB.GetUser(show_name)
  if err != nil {
    return nil, err
  } else {
    if PassHash != user.PassHash {
      return nil, fmt.Errorf("ERROR: Wrong password\n")
    } else {
      user.IsLogged = true
      DB.UpdateUser(user, user.ShowName)
      log.Printf("SUCCESS: User(%s) logged in", show_name)
      return user, nil
    }
  }
}

// User deactivation
func (DB *Database) LogOut(show_name string) {
  user, _ := DB.GetUser(show_name)
  user.IsLogged = false
  DB.UpdateUser(user, user.ShowName)
  log.Printf("SUCCESS: User(%s) logged out", show_name)
}

// Delete a user data
func (DB *Database) DelUser(show_name string) {
  user, err := DB.GetUser(show_name)
  if err != nil {
    log.Fatal("ERROR: ", err, "\n")
  } else {
    db := DB.DBConnect()
    defer db.Close()

    execSQL := `DELETE FROM users WHERE showname = $1`

    db.Exec(execSQL, user.ShowName)
    log.Printf("SUCCESS: User(%s) deleted", show_name)
  }
}

// List rooms
func (DB *Database) ListRooms() ([]string, []string) {
  db := DB.DBConnect()
  defer db.Close()

  execSQL := `
    SELECT name, hash FROM rooms
  `

  rows, err := db.Query(execSQL)
  defer rows.Close()

  var names []string
  var hashes []string
  for rows.Next() {
    var name string
    var hash string
    if err := rows.Scan(&name, &hash); err != nil {
      log.Println("ERROR: ", err)
      return nil, nil
    }
    names = append(names, name)
    hashes = append(hashes, hash)
  }

  if err != nil {
    log.Printf("ERROR: %s", err)
    return nil, nil
  }
  return names, hashes
}

// Getting room info
func (DB *Database) GetRoom(name string) (*Room, error) {
  db := DB.DBConnect()
  defer db.Close()

  execSQL := `
    SELECT name, creator_hash, hash, is_public, clients FROM rooms WHERE name = $1;
  `
  row := db.QueryRow(execSQL, name)
  var room Room
  err := row.Scan(
    &room.Name, &room.CreatorHash, &room.Hash,
    &room.IsPublic, pq.Array(&room.Clients),
  )
  if err != nil {
    log.Printf("ERROR: %s", err)
    return nil, err
  }

  return &room, nil
}

// Creates new room
func (DB *Database) BuildRoom(name string, creator_hash string, is_public bool) {
  db := DB.DBConnect()
  defer db.Close()

  roomHash := universal.CreateHash(map[string]interface{} {
    "Name":         name,
    "CreatorHash":  creator_hash,
  })

  room := Room {
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

  _, err := db.Exec(execSQL, &room.Name, &room.CreatorHash,
                    &room.Hash, &room.IsPublic, pq.Array(&room.Clients),
  )
  if err != nil {
    log.Printf("ERROR: %s", err)
  } else {
    log.Printf("SUCCESS: Room(%s) built", name)
  }
}

// Switches room privacy
func (DB *Database) SwitchRoomPriv(name string) {
  db := DB.DBConnect()
  defer db.Close()

  room, _ := DB.GetRoom(name)
  if room == nil {
    log.Printf("ERROR: Room(%s) doesn't exist", name)
    return
  }

  is_public := room.IsPublic
  if is_public {
    execSQL := `UPDATE rooms SET is_public = false WHERE name = $1`
    db.Exec(execSQL, name)
  } else {
    execSQL := `UPDATE rooms SET is_public = true WHERE name = $1`
    db.Exec(execSQL, name)
  }
}

// Close room
func (DB *Database) CloseRoom(name string) {
  db := DB.DBConnect()
  defer db.Close()

  _, err := DB.GetRoom(name)
  if err != nil {
    log.Printf("ERROR: Room(%s) doesn't exist", name)
    return
  }

  execSQL := `DELETE FROM rooms WHERE name = $1`
  db.Exec(execSQL, name)
}

// Updating room
func (DB *Database) UpdateRoom(new_room *Room, name string) {
  db := DB.DBConnect()
  defer db.Close()

	execSQL := `UPDATE rooms
              SET name = $1, creator_hash = $2, hash = $3, 
              is_public = $4, clients = $5
              WHERE name = $6`

  db.Exec(execSQL, new_room.Name, 
          new_room.CreatorHash, new_room.Hash, 
          new_room.IsPublic, pq.Array(new_room.Clients))

  log.Printf("SUCCESS: Room(%s) updated", new_room.Name)
}
