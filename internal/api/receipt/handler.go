package receipt

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type PointsResponse struct {
	Points int64 `json:"points"`
}

func (h *Handler) Points(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	points := h.service.Points(id)

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(PointsResponse{points})
}

type ProcessResponse struct {
	ID string `json:"id"`
}

func (h *Handler) Process(w http.ResponseWriter, r *http.Request) {
	id := h.service.SetPoints()

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(ProcessResponse{id})
}
