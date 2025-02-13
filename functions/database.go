package functions


import (
  "fmt"
  "log"
  "database/sql"
  "TermoChat/models"
  "TermoChat/config"
  _ "github.com/lib/pq"
)

// Converting user inputs to actual object
func UInput2User(input map[string]interface{}) *models.User {
  input["PassHash"] = CreateHash(map[string]interface{} {
    "Password": input["Password"].(string),
  })
  delete(input, "Password")
  input["Hash"] = CreateHash(input)

  user := models.User {
    ShowName:         input["ShowName"].(string),
    PassHash:         input["PassHash"].(string),
    RelatedQuestion:  input["RelatedQuestion"].(string),
    RelatedAnswer:    input["RelatedAnswer"].(string),
    Hash:             input["Hash"].(string),
    IsLogged:         input["IsLogged"].(bool),
  }
  return &user
}

// Simple db connection function
func DBConnect() *sql.DB{
  env := config.LoadEnv()
  connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", env.DB_USER, env.DB_PASS, env.DB_NAME)
  db, err := sql.Open("postgres", connStr)
  if err != nil {
    log.Fatal(err)
  }
  return db
}

// Database migration
func Migration() {
  db := DBConnect()
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

// Adding one user to database
func SignUp(user *models.User) {
  db := DBConnect()
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

// Searching for a user
func GetUser(ShowName string) (*models.User, error) {
  db := DBConnect()
  defer db.Close()

  execSQL := `SELECT show_name, pass_hash, related_question, related_answer, hash, is_logged
              FROM users WHERE show_name = $1`
  row := db.QueryRow(execSQL, ShowName)

  var user models.User
  err := row.Scan(&user.ShowName, &user.PassHash, &user.RelatedQuestion, &user.RelatedAnswer, &user.Hash, &user.IsLogged)
	if err != nil {
		if err == sql.ErrNoRows {
      return nil, fmt.Errorf("ERROR: User(%s) not found", ShowName)
		}
		return nil, err
	}
  return &user, nil
}

// Updating user on database
func UpdateUser(new_user *models.User, ShowName string) {
  db := DBConnect()
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
func LogIn(ShowName string, Password string) (*models.User, error){
  PassHash := CreateHash(map[string]interface{} {
    "Password": Password,
  })

  user, err := GetUser(ShowName)
  if err != nil {
    return nil, err
  } else {
    if PassHash != user.PassHash {
      return nil, fmt.Errorf("ERROR: Wrong password\n")
    } else {
      user.IsLogged = true
      UpdateUser(user, user.ShowName)
      log.Printf("ERROR: User(%s) logged in", ShowName)
      return user, nil
    }
  }
}

// User deactivation
func LogOut(ShowName string) {
  user, _ := GetUser(ShowName)
  user.IsLogged = false
  UpdateUser(user, user.ShowName)
  log.Printf("SUCCESS: User(%s) logged out", ShowName)
}

// Delete a user data
func DelUser(ShowName string) {
  user, err := GetUser(ShowName)
  if err != nil {
    log.Fatal("ERROR: ", err, "\n")
  } else {
    db := DBConnect()
    defer db.Close()

    execSQL := `DELETE FROM users WHERE showname = $1`

    db.Exec(execSQL, user.ShowName)
    log.Printf("SUCCESS: User(%s) deleted", ShowName)
  }
}

// List Public rooms
func ListPubRooms() []string {
  db := DBConnect()
  defer db.Close()

  execSQL := `
    SELECT hash FROM rooms
  `
  var hashes []string
  row := db.QueryRow(execSQL)
  err := row.Scan(&hashes)
  if err != nil {
    log.Println("ERROR: ", err)
    return nil
  } else {
    return hashes
  }
}
