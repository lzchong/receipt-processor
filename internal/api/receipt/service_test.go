package receipt

import (
	"errors"
	"testing"
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

func TestReceiptService_Points(t *testing.T) {
	repository := &stubRepository{}
	service := NewService(repository)

	t.Run("success", func(t *testing.T) {
		got, err := service.Points("7fb1377b-b223-49d9-a31a-5a02701dd310")
		assertEqual(t, "points", got, int64(32))
		assertNoError(t, err)
	})

	t.Run("zero points", func(t *testing.T) {
		got, err := service.Points("adb6b560-0eef-42bc-9d16-df48f30e89b2")
		assertEqual(t, "points", got, int64(0))
		assertNoError(t, err)
	})

	t.Run("missing ID", func(t *testing.T) {
		got, err := service.Points("6fb1377b-b223-49d9-a31a-5a02701dd310")
		assertEqual(t, "points", got, int64(0))
		assertError(t, err, ErrReceiptNotFound)
	})

	t.Run("empty ID", func(t *testing.T) {
		got, err := service.Points("")
		assertEqual(t, "points", got, int64(0))
		assertError(t, err, ErrReceiptNotFound)
	})
}

func TestReceiptService_SetPoints(t *testing.T) {
	mockRepo := &stubRepository{}
	service := NewService(mockRepo)

	t.Run("set points", func(t *testing.T) {
		assertEqual(t, "ID", service.SetPoints(), "7fb1377b-b223-49d9-a31a-5a02701dd310")
	})
}

func assertEqual[T comparable](t *testing.T, name string, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("expected %s %v, but got %v", name, want, got)
	}
}

func assertError(t *testing.T, got, want error) {
	t.Helper()
	if !errors.Is(got, want) {
		t.Errorf("expected error %v, but got %v", want, got)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
}
