package main

import (
	"log"
	"net/http"
)

func main() {
	if err := http.ListenAndServe(":8080", http.FileServer(http.Dir("./static"))); err != nil {
		log.Fatal(err)
	}
}
