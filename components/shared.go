package components


import (
  "sync"
)


var mu sync.Mutex
var db Database
