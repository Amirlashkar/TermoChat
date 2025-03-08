package universal


import (
  "crypto/sha256"
  "encoding/hex"
  "encoding/json"

  // "github.com/mitchellh/mapstructure"
	// "github.com/gorilla/websocket"
)


// Creating hash using secret key and data
func Data2Hash(data map[string]any) string {
  data["SECRET"] = secretKey
  jsonBytes, _ := json.Marshal(data)
  hash := sha256.Sum256(jsonBytes)
  return hex.EncodeToString(hash[:])
}

// To create hash from one word
func Word2Hash(word string) string {

  data := map[string]any {
    "WORD": word,
    "SECRET": secretKey,
  }

  bytes, _ := json.Marshal(data)
  hash := sha256.Sum256(bytes)
  return hex.EncodeToString(hash[:])
}
