package receipt

import (
	"errors"
	"testing"
	"time"
)

type stubRepository struct{}

func (m *stubRepository) Points(id string) (int64, bool) {
	switch id {
	case "7fb1377b-b223-49d9-a31a-5a02701dd310":
		return 32, true
	case "adb6b560-0eef-42bc-9d16-df48f30e89b2":
		return 0, true
	default:
		return 0, false
	}
}

func (m *stubRepository) CreatePoints(points int64) string {
	return "7fb1377b-b223-49d9-a31a-5a02701dd310"
}

func TestReceiptService_Points(t *testing.T) {
	repository := &stubRepository{}
	service := NewService(repository)

	tests := map[string]struct {
		input       string
		expected    int64
		expectedErr error
	}{
		"success":     {"7fb1377b-b223-49d9-a31a-5a02701dd310", 32, nil},
		"zero points": {"adb6b560-0eef-42bc-9d16-df48f30e89b2", 0, nil},
		"missing ID":  {"6fb1377b-b223-49d9-a31a-5a02701dd310", 0, ErrReceiptNotFound},
		"empty ID":    {"", 0, ErrReceiptNotFound},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := service.Points(test.input)
			if !errors.Is(err, test.expectedErr) {
				t.Errorf("expected error %v, but got %v", test.expectedErr, err)
			}
			want := test.expected
			if got != want {
				t.Errorf("expected points %d, but got %d", want, got)
			}
		})
	}
}

func TestReceiptService_Process(t *testing.T) {
	mockRepo := &stubRepository{}
	service := NewService(mockRepo)

	t.Run("set points", func(t *testing.T) {
		receiptItems := []ReceiptItem{{"", 0}}
		receipt := &Receipt{
			Retailer:     "M&M 0-1",
			PurchaseTime: time.Date(2024, time.December, 31, 13, 01, 0, 0, time.UTC),
			Items:        receiptItems,
			Total:        0,
		}
		if got, want := service.Process(receipt), "7fb1377b-b223-49d9-a31a-5a02701dd310"; got != want {
			t.Errorf("expected ID %v, but got %v", want, got)
		}
	})
}
