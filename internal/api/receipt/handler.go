package receipt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
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

var priceRegex = regexp.MustCompile(`^\d+\.\d{2}$`)
var descriptionRegex = regexp.MustCompile(`^[\w\s\-]+$`)
var retailerRegex = regexp.MustCompile(`^[\w\s\-&]+$`)

type ItemRequest struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

func (r *ItemRequest) Validate() error {
	if r.ShortDescription == "" {
		return fmt.Errorf("short description is required")
	}
	if matched := descriptionRegex.MatchString(r.ShortDescription); !matched {
		return fmt.Errorf("short description must contain only alphanumeric characters, spaces, and hyphens")
	}

	if !priceRegex.MatchString(r.Price) {
		return fmt.Errorf("price must be a decimal number with two decimal places")
	}

	return nil
}

func (r *ItemRequest) ToReceiptItem() (*ReceiptItem, error) {
	price, err := strconv.ParseFloat(r.Price, 64)
	if err != nil {
		return nil, err
	}
	item := &ReceiptItem{
		ShortDescription: r.ShortDescription,
		Price:            price,
	}
	return item, nil
}

type ProcessRequest struct {
	Retailer     string        `json:"retailer"`
	PurchaseDate string        `json:"purchaseDate"`
	PurchaseTime string        `json:"purchaseTime"`
	Items        []ItemRequest `json:"items"`
	Total        string        `json:"total"`
}

func (r *ProcessRequest) Validate() error {
	if r.Retailer == "" {
		return fmt.Errorf("retailer is required")
	}
	if matched := retailerRegex.MatchString(r.Retailer); !matched {
		return fmt.Errorf("retailer must contain only alphanumeric characters, spaces, hyphens, and ampersands")
	}

	if r.PurchaseDate == "" {
		return fmt.Errorf("purchase date is required")
	}
	if _, err := time.Parse(time.DateOnly, r.PurchaseDate); err != nil {
		return fmt.Errorf("purchase date is not a valid date, %v", err)
	}

	if r.PurchaseTime == "" {
		return fmt.Errorf("purchase time is required")
	}
	if _, err := time.Parse("15:04", r.PurchaseTime); err != nil {
		return fmt.Errorf("purchase time is not a valid time, %v", err)
	}

	if len(r.Items) == 0 {
		return fmt.Errorf("minimum of one item is required")
	}
	for _, item := range r.Items {
		if err := item.Validate(); err != nil {
			return fmt.Errorf("item %s is invalid: %v", item.ShortDescription, err)
		}
	}

	if !priceRegex.MatchString(r.Total) {
		return fmt.Errorf("total must be a decimal number with two decimal places")
	}

	return nil
}

func (r *ProcessRequest) ToReceipt() (*Receipt, error) {
	purchaseTime, err := time.Parse(time.DateTime, fmt.Sprintf("%s %s:00", r.PurchaseDate, r.PurchaseTime))
	if err != nil {
		return nil, err
	}

	items := make([]ReceiptItem, len(r.Items))
	for i, itemRequest := range r.Items {
		item, err := itemRequest.ToReceiptItem()
		if err != nil {
			return nil, err
		}
		items[i] = *item
	}

	total, err := strconv.ParseFloat(r.Total, 64)
	if err != nil {
		return nil, err
	}
	receipt := &Receipt{
		Retailer:     r.Retailer,
		PurchaseTime: purchaseTime,
		Items:        items,
		Total:        total,
	}

	return receipt, nil
}

type ProcessResponse struct {
	ID string `json:"id"`
}

func (h *handlerImpl) Process(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Missing request body. Please provide a JSON object representing a receipt.", http.StatusBadRequest)
		return
	}

	maxBodySize := int64(1 << 20) // 1MB limit
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	var dto ProcessRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&dto); err != nil {
		http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
		return
	}

	if err := dto.Validate(); err != nil {
		http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
		return
	}

	receipt, err := dto.ToReceipt()
	if err != nil {
		http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
		return
	}

	id := h.service.Process(receipt)

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(ProcessResponse{id})
}
