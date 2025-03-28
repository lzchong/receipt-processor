package receipt

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type stubService struct{}

func (m *stubService) Points(id string) (int64, error) {
	if id == "7fb1377b-b223-49d9-a31a-5a02701dd310" {
		return 32, nil
	}
	return 0, ErrReceiptNotFound
}

func (m *stubService) Process(receipt *Receipt) string {
	return "7fb1377b-b223-49d9-a31a-5a02701dd310"
}

func TestReceiptHandler_Points(t *testing.T) {
	service := &stubService{}
	handler := NewHandler(service)
	mux := http.NewServeMux()
	mux.HandleFunc("/receipts/{id}/points", handler.Points)

	t.Run("success", func(t *testing.T) {
		request, err := http.NewRequest("GET", "/receipts/7fb1377b-b223-49d9-a31a-5a02701dd310/points", nil)
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
		assertContentType(t, response, "application/json")
		assertJSONResponse(t, response, PointsResponse{32})
	})

	t.Run("affix whilespace", func(t *testing.T) {
		request, err := http.NewRequest("GET", "/receipts/%207fb1377b-b223-49d9-a31a-5a02701dd310%20/points", nil)
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
		assertContentType(t, response, "application/json")
		assertJSONResponse(t, response, PointsResponse{32})
	})

	t.Run("invalid ID", func(t *testing.T) {
		request, err := http.NewRequest("GET", "/receipts/7fb1377b-b223-49d9%20a31a-5a02701dd310/points", nil)
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusBadRequest)
		assertHasError(t, response)
	})

	t.Run("not found", func(t *testing.T) {
		request, err := http.NewRequest("GET", "/receipts/non-existent-id/points", nil)
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusNotFound)
		assertHasError(t, response)
	})

	t.Run("empty ID", func(t *testing.T) {
		request, err := http.NewRequest("GET", "/receipts//points", nil)
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusMovedPermanently)
	})
}

func TestReceiptHandler_ItemRequestValidation(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		receiptItem := &ItemRequest{"Mountain Dew 12PK", "6.49"}
		err := receiptItem.Validate()
		if err != nil {
			t.Errorf("expected no error, but got %v", err)
		}
	})

	t.Run("empty short description", func(t *testing.T) {
		receiptItem := &ItemRequest{"", "6.49"}
		err := receiptItem.Validate()
		if err == nil {
			t.Error("expected has error, but got nothing")
		}
	})

	t.Run("invalid short description", func(t *testing.T) {
		receiptItem := &ItemRequest{"Mountain&Dew 12PK", "6.49"}
		err := receiptItem.Validate()
		if err == nil {
			t.Error("expected has error, but got nothing")
		}
	})

	t.Run("price lower than zero", func(t *testing.T) {
		receiptItem := &ItemRequest{"Mountain Dew 12PK", "-1"}
		err := receiptItem.Validate()
		if err == nil {
			t.Error("expected has error, but got nothing")
		}
	})

	t.Run("price in incorrect decimal places", func(t *testing.T) {
		receiptItem := &ItemRequest{"Mountain Dew 12PK", "6.496"}
		err := receiptItem.Validate()
		if err == nil {
			t.Error("expected has error, but got nothing")
		}
	})
}

func TestReceiptHandler_ProcessRequestValidation(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		receipt := &ProcessRequest{
			Retailer:     "Target",
			PurchaseDate: "2022-01-01",
			PurchaseTime: "13:01",
			Items:        []ItemRequest{{"Mountain Dew 12PK", "6.49"}},
			Total:        "35.35",
		}
		err := receipt.Validate()
		if err != nil {
			t.Errorf("expected no error, but got %v", err)
		}
	})

	t.Run("empty retailer", func(t *testing.T) {
		receipt := &ProcessRequest{
			Retailer:     "",
			PurchaseDate: "2022-01-01",
			PurchaseTime: "13:01",
			Items:        []ItemRequest{{"Mountain Dew 12PK", "6.49"}},
			Total:        "35.35",
		}
		err := receipt.Validate()
		if err == nil {
			t.Error("expected has error, but got nothing")
		}
	})

	t.Run("empty purchase date", func(t *testing.T) {
		receipt := &ProcessRequest{
			Retailer:     "Target",
			PurchaseDate: "",
			PurchaseTime: "13:01",
			Items:        []ItemRequest{{"Mountain Dew 12PK", "6.49"}},
			Total:        "35.35",
		}
		err := receipt.Validate()
		if err == nil {
			t.Error("expected has error, but got nothing")
		}
	})

	t.Run("invalid purchase date", func(t *testing.T) {
		receipt := &ProcessRequest{
			Retailer:     "Target",
			PurchaseDate: "01-01-2022",
			PurchaseTime: "13:01",
			Items:        []ItemRequest{{"Mountain Dew 12PK", "6.49"}},
			Total:        "35.35",
		}
		err := receipt.Validate()
		if err == nil {
			t.Error("expected has error, but got nothing")
		}
	})

	t.Run("empty purchase time", func(t *testing.T) {
		receipt := &ProcessRequest{
			Retailer:     "Target",
			PurchaseDate: "2022-01-01",
			PurchaseTime: "",
			Items:        []ItemRequest{{"Mountain Dew 12PK", "6.49"}},
			Total:        "35.35",
		}
		err := receipt.Validate()
		if err == nil {
			t.Error("expected has error, but got nothing")
		}
	})

	t.Run("invalid purchase time", func(t *testing.T) {
		receipt := &ProcessRequest{
			Retailer:     "Target",
			PurchaseDate: "2022-01-01",
			PurchaseTime: "13:01:32",
			Items:        []ItemRequest{{"Mountain Dew 12PK", "6.49"}},
			Total:        "35.35",
		}
		err := receipt.Validate()
		if err == nil {
			t.Error("expected has error, but got nothing")
		}
	})

	t.Run("contains invalid item", func(t *testing.T) {
		receipt := &ProcessRequest{
			Retailer:     "Target",
			PurchaseDate: "2022-01-01",
			PurchaseTime: "13:01",
			Items:        []ItemRequest{{"", "6.49"}},
			Total:        "35.35",
		}
		err := receipt.Validate()
		if err == nil {
			t.Error("expected has error, but got nothing")
		}
	})

	t.Run("no items", func(t *testing.T) {
		receipt := &ProcessRequest{
			Retailer:     "Target",
			PurchaseDate: "2022-01-01",
			PurchaseTime: "13:01",
			Items:        []ItemRequest{},
			Total:        "35.35",
		}
		err := receipt.Validate()
		if err == nil {
			t.Error("expected has error, but got nothing")
		}
	})

	t.Run("total lower than zero", func(t *testing.T) {
		receipt := &ProcessRequest{
			Retailer:     "Target",
			PurchaseDate: "2022-01-01",
			PurchaseTime: "13:01",
			Items:        []ItemRequest{{"Mountain Dew 12PK", "6.49"}},
			Total:        "-1",
		}
		err := receipt.Validate()
		if err == nil {
			t.Error("expected has error, but got nothing")
		}
	})

	t.Run("total in incorrect decimal places", func(t *testing.T) {
		receipt := &ProcessRequest{
			Retailer:     "Target",
			PurchaseDate: "2022-01-01",
			PurchaseTime: "13:01",
			Items:        []ItemRequest{{"Mountain Dew 12PK", "6.49"}},
			Total:        "35.353",
		}
		err := receipt.Validate()
		if err == nil {
			t.Error("expected has error, but got nothing")
		}
	})
}

func TestReceiptHandler_Process(t *testing.T) {
	service := &stubService{}
	handler := NewHandler(service)

	t.Run("success", func(t *testing.T) {
		receipt := &ProcessRequest{
			Retailer:     "Target",
			PurchaseDate: "2022-01-01",
			PurchaseTime: "13:01",
			Items: []ItemRequest{
				{"Mountain Dew 12PK", "6.49"},
				{"Emils Cheese Pizza", "12.25"},
				{"Knorr Creamy Chicken", "1.26"},
				{"Doritos Nacho Cheese", "3.35"},
				{"   Klarbrunn 12-PK 12 FL OZ  ", "12.00"},
			},
			Total: "35.35",
		}

		body, err := json.Marshal(receipt)
		if err != nil {
			t.Fatalf("failed to marshal receipt, %v", err)
		}

		request, err := http.NewRequest("POST", "/receipts/process", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()
		handler.Process(response, request)

		assertStatus(t, response, http.StatusAccepted)
		assertContentType(t, response, "application/json")
		assertJSONResponse(t, response, ProcessResponse{"7fb1377b-b223-49d9-a31a-5a02701dd310"})
	})

	t.Run("no request body", func(t *testing.T) {
		request, err := http.NewRequest("POST", "/receipts/process", nil)
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()
		handler.Process(response, request)

		assertStatus(t, response, http.StatusBadRequest)
		assertHasError(t, response)
	})

	t.Run("malformed JSON", func(t *testing.T) {
		body := `{retailer":"Target","purchaseDate":"2022-01-01","purchaseTime":"13:01","items":[{"shortDescription":"Mountain Dew 12PK","price":"6.49"}],"total":"6.49"}`
		request, err := http.NewRequest("POST", "/receipts/process", strings.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()
		handler.Process(response, request)

		assertStatus(t, response, http.StatusBadRequest)
		assertHasError(t, response)
	})

	t.Run("unexpected fields", func(t *testing.T) {
		body := `{"extra-field":"some value","retailer":"Target","purchaseDate":"2022-01-01","purchaseTime":"13:01","items":[{"shortDescription":"Mountain Dew 12PK","price":"6.49"}],"total":"6.49"}`
		request, err := http.NewRequest("POST", "/receipts/process", strings.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()
		handler.Process(response, request)

		assertStatus(t, response, http.StatusBadRequest)
		assertHasError(t, response)
	})

	t.Run("invalid data types", func(t *testing.T) {
		body := `{"retailer":"Target","purchaseDate":"2022-01-01","purchaseTime":"13:01","items":[{"shortDescription":"Mountain Dew 12PK","price":"6.49"}],"total":"price"}`
		request, err := http.NewRequest("POST", "/receipts/process", strings.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()
		handler.Process(response, request)

		assertStatus(t, response, http.StatusBadRequest)
		assertHasError(t, response)
	})
}

func assertStatus(t *testing.T, response *httptest.ResponseRecorder, want int) {
	t.Helper()
	if got := response.Code; got != want {
		t.Errorf("expected status %d, but got %d", want, got)
	}
}

func assertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if got := response.Result().Header.Get("content-type"); got != want {
		t.Errorf("expected content-type %s, but got %s", want, got)
	}
}

func assertHasError(t *testing.T, response *httptest.ResponseRecorder) {
	t.Helper()
	if got := strings.TrimSpace(response.Body.String()); got == "" {
		t.Error("expected error message, but got nothing")
	}
}

func assertJSONResponse[T any](t *testing.T, response *httptest.ResponseRecorder, want T) {
	t.Helper()

	var got T
	decoder := json.NewDecoder(response.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&got); err != nil {
		t.Fatalf("failed to parse response %q, '%v'", response.Body, err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected body %v, but got %v", want, got)
	}
}
