package main

import (
  "fmt"
  "net/http"

  "TermoChat/routers"
)

func main() {
    // Gets general router not to get lost :)
    router := routers.ProvideRouter()
    fmt.Println("Server going live on :8000 ...")
    http.ListenAndServe(":8000", router)
}
