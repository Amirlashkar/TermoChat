package functions


import (
  "crypto/sha256"
  "encoding/hex"
  "encoding/json"
  "log"
  "fmt"
  "github.com/mitchellh/mapstructure"
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

// Creates a hash of some sample map data
func CreateDataHash(data map[string]interface{}) string {
  jsonBytes, err := json.Marshal(data)
  if err != nil {
    log.Fatalf("Error: %v", err)
  }

  hash := sha256.Sum256(jsonBytes)
  return hex.EncodeToString(hash[:])
}
