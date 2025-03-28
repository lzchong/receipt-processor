package receipt

import (
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

func (m *stubService) SetPoints() string {
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

func TestReceiptHandler_Process(t *testing.T) {
	service := &stubService{}
	handler := NewHandler(service)

	t.Run("process", func(t *testing.T) {
		request, err := http.NewRequest("POST", "/receipts/process", nil)
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()
		handler.Process(response, request)

		assertStatus(t, response, http.StatusAccepted)
		assertContentType(t, response, "application/json")
		assertJSONResponse(t, response, ProcessResponse{"7fb1377b-b223-49d9-a31a-5a02701dd310"})
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
