package functions


import (
  "fmt"
  "log"
  "crypto/sha256"
  "encoding/hex"
  "encoding/json"
  "github.com/mitchellh/mapstructure"
  "TermoChat/config"
)


// Makes a custom object from its data map
func Map2Obj(data map[string]interface{}, obj interface{}) interface{} {
  err := mapstructure.Decode(data, &obj)
  if err != nil {
    fmt.Printf("mapstructure Error: %s", err)
    return nil
  }
  return obj
}

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
