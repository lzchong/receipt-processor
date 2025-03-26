package receipt

import (
	"testing"
)

func TestReceiptService(t *testing.T) {
	service := NewService()

	t.Run("points", func(t *testing.T) {
		id := "7fb1377b-b223-49d9-a31a-5a02701dd310"
		assertEqual(t, "points", service.Points(id), int64(32))
	})

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
