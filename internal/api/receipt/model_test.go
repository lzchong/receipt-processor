package receipt

import (
	"testing"
	"time"
)

func TestCalculatePoints(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		receipt := &Receipt{
			Retailer:     "Target",
			PurchaseTime: time.Date(2022, time.January, 1, 13, 01, 0, 0, time.UTC),
			Items: []ReceiptItem{
				{"Mountain Dew 12PK", 6.49},
				{"Emils Cheese Pizza", 12.25},
				{"Knorr Creamy Chicken", 1.26},
				{"Doritos Nacho Cheese", 3.35},
				{"   Klarbrunn 12-PK 12 FL OZ  ", 12.00},
			},
			Total: 35.35,
		}

		got := receipt.CalculatePoints()
		assertEqual(t, "points", got, 28)
	})
}

func TestCountByAlphanumericCharacter(t *testing.T) {
	t.Run("return 14 for M&M Corner Market", func(t *testing.T) {
		got := countByAlphanumericCharacter("M&M Corner Market")
		assertEqual(t, "count", got, 14)
	})

	t.Run("return 0 for empty string", func(t *testing.T) {
		got := countByAlphanumericCharacter("")
		assertEqual(t, "count", got, 0)
	})
}

func TestIsRoundDollarAmount(t *testing.T) {
	t.Run("return true for 0.00", func(t *testing.T) {
		got := isRoundDollarAmount(0.00)
		assertTrue(t, got)
	})

	t.Run("return false for 0.01", func(t *testing.T) {
		got := isRoundDollarAmount(0.01)
		assertFalse(t, got)
	})

	t.Run("return false for 0.99", func(t *testing.T) {
		got := isRoundDollarAmount(0.99)
		assertFalse(t, got)
	})

	t.Run("return false for 1.00", func(t *testing.T) {
		got := isRoundDollarAmount(1.00)
		assertTrue(t, got)
	})
}

func TestIsMultipleOfQuarter(t *testing.T) {
	t.Run("return true for 0.00", func(t *testing.T) {
		got := isMultipleOfQuarter(0.00)
		assertTrue(t, got)
	})

	t.Run("return true for 5.75", func(t *testing.T) {
		got := isMultipleOfQuarter(5.75)
		assertTrue(t, got)
	})

	t.Run("return false for 25.01", func(t *testing.T) {
		got := isMultipleOfQuarter(25.01)
		assertFalse(t, got)
	})
}

func TestCountEveryTwoItems(t *testing.T) {
	t.Run("return 1 for two items", func(t *testing.T) {
		got := countEveryTwoItems([]ReceiptItem{{"", 0}, {"", 0}})
		assertEqual(t, "count", got, 1)
	})

	t.Run("return 0 for one item", func(t *testing.T) {
		got := countEveryTwoItems([]ReceiptItem{{"", 0}})
		assertEqual(t, "count", got, 0)
	})

	t.Run("return 0 for no item", func(t *testing.T) {
		got := countEveryTwoItems([]ReceiptItem{})
		assertEqual(t, "count", got, 0)
	})
}

func TestIsStringLengthMultipleOfThree(t *testing.T) {
	t.Run("return true for 3 characters", func(t *testing.T) {
		got := isStringLengthMultipleOfThree("abc")
		assertTrue(t, got)
	})

	t.Run("return false for 4 characters", func(t *testing.T) {
		got := isStringLengthMultipleOfThree("abcd")
		assertFalse(t, got)
	})

	t.Run("return true for 0 character", func(t *testing.T) {
		got := isStringLengthMultipleOfThree("")
		assertTrue(t, got)
	})
}

func TestIsOddDay(t *testing.T) {
	t.Run("return true if day is 1", func(t *testing.T) {
		d, _ := time.Parse(time.DateOnly, "2024-01-01")
		got := isOddDay(d)
		assertTrue(t, got)
	})

	t.Run("return false if day is 2", func(t *testing.T) {
		d, _ := time.Parse(time.DateOnly, "2024-01-02")
		got := isOddDay(d)
		assertFalse(t, got)
	})
}

func TestIsTimeBetweenTwoPMAndFourPM(t *testing.T) {
	t.Run("return true if time is 2:00PM", func(t *testing.T) {
		ti, _ := time.Parse("15:04", "14:00")
		got := isTimeBetweenTwoPMAndFourPM(ti)
		assertFalse(t, got)
	})

	t.Run("return true if time is 2:01PM", func(t *testing.T) {
		ti, _ := time.Parse("15:04", "14:01")
		got := isTimeBetweenTwoPMAndFourPM(ti)
		assertTrue(t, got)
	})

	t.Run("return true if time is 3:59PM", func(t *testing.T) {
		ti, _ := time.Parse("15:04", "15:59")
		got := isTimeBetweenTwoPMAndFourPM(ti)
		assertTrue(t, got)
	})

	t.Run("return false if time is 4:00PM", func(t *testing.T) {
		ti, _ := time.Parse("15:04", "16:00")
		got := isTimeBetweenTwoPMAndFourPM(ti)
		assertFalse(t, got)
	})
}

func assertTrue(t *testing.T, got bool) {
	t.Helper()
	if got != true {
		t.Errorf("expected true, but got %t", got)
	}
}

func assertFalse(t *testing.T, got bool) {
	t.Helper()
	if got != false {
		t.Errorf("expected false, but got %t", got)
	}
}
