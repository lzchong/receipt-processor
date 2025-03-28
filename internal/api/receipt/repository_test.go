package receipt

import (
	"testing"
)

func TestReceiptRepository(t *testing.T) {
	repo := NewRepository()

	t.Run("get points", func(t *testing.T) {
		id := "test-id"
		points := int64(100)

		repo.(*inMemoryRepository).points[id] = points

		got, ok := repo.Points(id)
		assertPoints(t, got, ok, points)
	})

	t.Run("not found", func(t *testing.T) {
		id := "non-existent-id"

		got, ok := repo.Points(id)
		assertNoPoints(t, got, ok)
	})

	t.Run("create points", func(t *testing.T) {
		points := int64(100)

		id := repo.CreatePoints(points)
		if id == "" {
			t.Fatal("expected ID, but got nothing")
		}

		got := repo.(*inMemoryRepository).points[id]
		assertPoints(t, got, true, points)
	})
}

func assertPoints(t *testing.T, got int64, ok bool, want int64) {
	t.Helper()
	if !ok {
		t.Fatalf("expected points %d, but got nil", want)
	}
	if got != want {
		t.Errorf("expected points %d, but got %d", want, got)
	}
}

func assertNoPoints(t *testing.T, got int64, ok bool) {
	t.Helper()
	if ok {
		t.Errorf("expected no points, but got %d", got)
	}
	if got != 0 {
		t.Errorf("expected points 0, but got %d", got)
	}
}
