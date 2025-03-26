package server

import (
	"github.com/lzchong/receipt-processor/internal/api/receipt"
	"net/http"
)

func NewRouter() http.Handler {
	receiptService := receipt.NewService()
	receiptHandler := receipt.NewHandler(receiptService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /receipts/{id}/points", receiptHandler.Points)
	mux.HandleFunc("POST /receipts/process", receiptHandler.Process)

	return mux
}
