package main

import (
	"github.com/lzchong/receipt-processor/internal/server"
	"log"
)

func main() {
	router := server.NewRouter()
	s := server.NewServer(router)
	log.Println("Starting server on :8080...")
	log.Fatal(s.ListenAndServe())
}
