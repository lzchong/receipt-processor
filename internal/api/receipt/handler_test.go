package receipt

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type mockService struct{}

func (m *mockService) Points(id string) int64 {
	return 32
}

func (m *mockService) SetPoints() string {
	return "7fb1377b-b223-49d9-a31a-5a02701dd310"
}

func TestReceiptHandler(t *testing.T) {
	service := &mockService{}
	handler := NewHandler(service)

	t.Run("points", func(t *testing.T) {
		request, err := http.NewRequest("GET", "/receipts/7fb1377b-b223-49d9-a31a-5a02701dd310/points", nil)
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()
		handler.Points(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, "application/json")
		assertJSONResponse(t, response, PointsResponse{32})
	})

	t.Run("process", func(t *testing.T) {
		request, err := http.NewRequest("POST", "/receipts/process", nil)
		if err != nil {
			t.Fatal(err)
		}

		response := httptest.NewRecorder()
		handler.Process(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)
		assertContentType(t, response, "application/json")
		assertJSONResponse(t, response, ProcessResponse{"7fb1377b-b223-49d9-a31a-5a02701dd310"})
	})
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("expected status %d, but got %d", want, got)
	}
}

func assertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if got := response.Result().Header.Get("content-type"); got != want {
		t.Errorf("expected content-type %s, but got %s", want, got)
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
