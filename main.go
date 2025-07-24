package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, World!")
    })
    http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
        name := r.FormValue("name")
        fmt.Fprintf(w, "Hello, %s!",name)
    })
    log.Fatal(http.ListenAndServe(":8080", nil))
}