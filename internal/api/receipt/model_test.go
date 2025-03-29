package receipt

import (
	"testing"
	"time"
)

func TestCalculatePoints(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		receipt := &Receipt{
			Retailer:     "Target",
			PurchaseTime: time.Date(2022, time.January, 1, 14, 01, 0, 0, time.UTC),
			Items: []ReceiptItem{
				{"Mountain Dew 12PK", 6.49},
				{"Emils Cheese Pizza", 12.25},
				{"Knorr Creamy Chicken", 1.26},
				{"Doritos Nacho Cheese", 3.35},
				{"   Klarbrunn 12-PK 12 FL OZ  ", 12.00},
			},
			Total: 35.35,
		}

		if got, want := receipt.CalculatePoints(), int64(38); got != want {
			t.Errorf("expected points %v, but got %v", want, got)
		}
	})
}

func TestCountByAlphanumericCharacter(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected int64
	}{
		"M&M Corner Market": {"M&M Corner Market", 14},
		"empty string":      {"", 0},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got, want := countByAlphanumericCharacter(test.input), test.expected; got != want {
				t.Errorf("expected count %v, but got %v", want, got)
			}
		})
	}
}

func TestIsRoundDollarAmount(t *testing.T) {
	tests := map[string]struct {
		input    float64
		expected bool
	}{
		"0.00": {0.00, true},
		"0.01": {0.01, false},
		"0.99": {0.99, false},
		"1.00": {1.00, true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got, want := isRoundDollarAmount(test.input), test.expected; got != want {
				t.Errorf("expected %t, but got %t", want, got)
			}
		})
	}
}

func TestIsMultipleOfQuarter(t *testing.T) {
	tests := map[string]struct {
		input    float64
		expected bool
	}{
		"0.00":  {0.00, true},
		"5.75":  {5.75, true},
		"25.01": {25.01, false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got, want := isMultipleOfQuarter(test.input), test.expected; got != want {
				t.Errorf("expected %t, but got %t", want, got)
			}
		})
	}
}

func TestCountEveryTwoItems(t *testing.T) {
	tests := map[string]struct {
		input    []ReceiptItem
		expected int64
	}{
		"two items": {[]ReceiptItem{{"", 0}, {"", 0}}, 1},
		"one item":  {[]ReceiptItem{{"", 0}}, 0},
		"no item":   {[]ReceiptItem{}, 0},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got, want := countEveryTwoItems(test.input), test.expected; got != want {
				t.Errorf("expected count %d, but got %d", want, got)
			}
		})
	}
}

func TestIsStringLengthMultipleOfThree(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected bool
	}{
		"three characters": {"abc", true},
		"four characters":  {"abcd", false},
		"no character":     {"", true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got, want := isStringLengthMultipleOfThree(test.input), test.expected; got != want {
				t.Errorf("expected %t, but got %t", want, got)
			}
		})
	}
}

func TestIsOddDay(t *testing.T) {
	tests := map[string]struct {
		input    int
		expected bool
	}{
		"day 1": {1, true},
		"day 2": {2, false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			inputTime := time.Date(2024, time.January, test.input, 0, 0, 0, 0, time.UTC)
			if got, want := isOddDay(inputTime), test.expected; got != want {
				t.Errorf("expected %t, but got %t", want, got)
			}
		})
	}
}

func TestIsTimeBetweenTwoPMAndFourPM(t *testing.T) {
	tests := map[string]struct {
		hour     int
		min      int
		expected bool
	}{
		"2:00 PM": {14, 0, false},
		"2:01 PM": {14, 1, true},
		"3:59 PM": {15, 59, true},
		"4:00 PM": {16, 0, false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			inputTime := time.Date(2024, time.January, 1, test.hour, test.min, 0, 0, time.UTC)
			if got, want := isTimeBetweenTwoPMAndFourPM(inputTime), test.expected; got != want {
				t.Errorf("expected %t, but got %t", want, got)
			}
		})
	}
}
