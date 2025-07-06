package main

import (
	"log"
	"net/http"
)

const addr = ":3000"

func main() {
	log.Printf("Look at http://localhost%v/", addr)
	// Используем FileServer для статических файлов
	fileServer := http.FileServer(http.Dir("static"))
	if err := http.ListenAndServe(addr, fileServer); err != nil { //nolint:gosec // non-prod solution
		log.Fatal(err)
	}
}
