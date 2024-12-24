package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World!")
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/", HelloWorld)
    http.Handle("/", r)

    fmt.Println("Server is running on port 8000...")
    log.Fatal(http.ListenAndServe(":8000", nil))
}



