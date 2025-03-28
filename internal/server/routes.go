package server

import (
	"github.com/lzchong/receipt-processor/internal/api/receipt"
	"net/http"
)

func NewRouter(receiptHandler receipt.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /receipts/{id}/points", receiptHandler.Points)
	mux.HandleFunc("POST /receipts/process", receiptHandler.Process)
	return mux
}
