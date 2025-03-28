package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type stubHandler struct{}

func (h *stubHandler) Points(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *stubHandler) Process(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
}

func TestRouter(t *testing.T) {
	handler := &stubHandler{}
	router := NewRouter(handler)

	testCases := map[string]struct {
		method string
		path   string
		want   int
	}{
		"get correct points":      {"GET", "/receipts/7fb1377b-b223-49d9-a31a-5a02701dd310/points", http.StatusOK},
		"trailing slash":          {"GET", "/receipts/7fb1377b-b223-49d9-a31a-5a02701dd310/points/", http.StatusNotFound},
		"missing ID":              {"GET", "/receipts//points", http.StatusMovedPermanently},
		"process receipt success": {"POST", "/receipts/process", http.StatusAccepted},
		"unsupported method":      {"DELETE", "/receipts/process", http.StatusMethodNotAllowed},
		"invalid path":            {"GET", "/invalid/route", http.StatusNotFound},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request, err := http.NewRequest(tc.method, tc.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			assertStatus(t, response, tc.want)
		})
	}
}

func assertStatus(t *testing.T, response *httptest.ResponseRecorder, want int) {
	t.Helper()
	if got := response.Code; got != want {
		t.Errorf("expected status %d, but got %d", want, got)
	}
}
