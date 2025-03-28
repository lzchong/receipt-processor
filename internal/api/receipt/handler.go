package receipt

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
)

type Handler interface {
	Points(w http.ResponseWriter, r *http.Request)
	Process(w http.ResponseWriter, r *http.Request)
}

type handlerImpl struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handlerImpl{service}
}

type PointsResponse struct {
	Points int64 `json:"points"`
}

var noWhitespaceRegex = regexp.MustCompile("^\\S+$")

func (h *handlerImpl) Points(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))

	if id == "" {
		http.Error(w, "Receipt ID cannot be empty.", http.StatusBadRequest)
		return
	}

	matched := noWhitespaceRegex.MatchString(id)
	if !matched {
		http.Error(w, "Receipt ID is invalid.", http.StatusBadRequest)
		return
	}

	points, err := h.service.Points(id)
	if err != nil {
		http.Error(w, "No receipt found for that ID.", http.StatusNotFound)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(PointsResponse{points})
}

type ProcessResponse struct {
	ID string `json:"id"`
}

func (h *handlerImpl) Process(w http.ResponseWriter, r *http.Request) {
	id := h.service.SetPoints()

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(ProcessResponse{id})
}
