package main

import (
	"github.com/lzchong/receipt-processor/internal/api/receipt"
	"github.com/lzchong/receipt-processor/internal/server"
	"log"
)

func main() {
	receiptRepository := receipt.NewRepository()
	receiptService := receipt.NewService(receiptRepository)
	receiptHandler := receipt.NewHandler(receiptService)

	router := server.NewRouter(receiptHandler)
	s := server.NewServer(router)

	log.Println("Starting server on :8080...")
	log.Fatal(s.ListenAndServe())
}
