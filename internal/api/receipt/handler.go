package receipt

import (
	"encoding/json"
	"net/http"
)

type IDResponse struct {
	ID string `json:"id"`
}

type PointsResponse struct {
	Points int64 `json:"points"`
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Points(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	points := h.service.Points(id)

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(PointsResponse{points})
}

func (h *Handler) SetPoints(w http.ResponseWriter, r *http.Request) {
	id := h.service.SetPoints()

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(IDResponse{id})
}
