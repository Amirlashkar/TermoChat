package config


import (
  "os"
  "log"
  "github.com/joho/godotenv"
)


type Env struct {
  SECRET          string
  PORT            string
  DB_NAME         string
  DB_USER         string
  DB_PASS         string
  SERVER_ADMIN    string
}

func LoadEnv() Env {
  err := godotenv.Load()
  if err != nil {
    log.Println("Warning: No .env file found, using system environment variables")
  }

  return Env {
    SECRET:       os.Getenv("SECRET_KEY"),
    PORT:         os.Getenv("PORT"),
    DB_NAME:      os.Getenv("DB_NAME"),
    DB_USER:      os.Getenv("DB_USER"),
    DB_PASS:      os.Getenv("DB_PASS"),
    SERVER_ADMIN: os.Getenv("SERVER_ADMIN"),
  }
}
