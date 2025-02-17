package universal


import (
  // "fmt"
  "log"
  "crypto/sha256"
  "encoding/hex"
  "encoding/json"

  "TermoChat/config"

  // "github.com/mitchellh/mapstructure"
	// "github.com/gorilla/websocket"
)

// Creating hash using secret key and data
func CreateHash(data map[string]interface{}) string {
  secretKey := config.LoadEnv().SECRET
  data["SECRET"] = secretKey
  jsonBytes, err := json.Marshal(data)
  if err != nil {
    log.Fatalf("Error: %v", err)
  }

  hash := sha256.Sum256(jsonBytes)
  return hex.EncodeToString(hash[:])
}


