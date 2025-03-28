package receipt

import (
	"math"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type ReceiptItem struct {
	ShortDescription string
	Price            float64
}

type Receipt struct {
	Retailer     string
	PurchaseTime time.Time
	Items        []ReceiptItem
	Total        float64
}

func (r *Receipt) CalculatePoints() int64 {
	points := int64(0)

	// One point for every alphanumeric character in the retailer
	points += 1 * countByAlphanumericCharacter(r.Retailer)

	// 50 points if the total is a round dollar amount with no cents
	if isRoundDollarAmount(r.Total) {
		points += 50
	}

	// 25 points if the total is a multiple of 0.25
	if isMultipleOfQuarter(r.Total) {
		points += 25
	}

	// 5 points for every two items on the receipt
	points += 5 * countEveryTwoItems(r.Items)

	// If the trimmed length of the item description is a multiple of 3,
	// multiply the price by 0.2 and round up to the nearest integer.
	// The result is the number of points earned.
	for _, item := range r.Items {
		if isStringLengthMultipleOfThree(item.ShortDescription) {
			points += int64(math.Ceil(item.Price * 0.2))
		}
	}

	// 6 points if the day in the purchase date is odd
	if isOddDay(r.PurchaseTime) {
		points += 6
	}

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm
	if isTimeBetweenTwoPMAndFourPM(r.PurchaseTime) {
		points += 10
	}

	return points
}

func countByAlphanumericCharacter(s string) int64 {
	count := int64(0)
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			count += 1
		}
	}
	return count
}

func isRoundDollarAmount(total float64) bool {
	return total == math.Round(total)
}

func isMultipleOfQuarter(total float64) bool {
	return total == math.Round(total*4)/4
}

func countEveryTwoItems(items []ReceiptItem) int64 {
	return int64(len(items) / 2)
}

func isStringLengthMultipleOfThree(s string) bool {
	return utf8.RuneCountInString(strings.TrimSpace(s))%3 == 0
}

func isOddDay(t time.Time) bool {
	return t.Day()%2 == 1
}

func isTimeBetweenTwoPMAndFourPM(t time.Time) bool {
	twoPM := time.Date(t.Year(), t.Month(), t.Day(), 14, 0, 0, 0, t.Location())
	fourPM := time.Date(t.Year(), t.Month(), t.Day(), 16, 0, 0, 0, t.Location())
	return t.After(twoPM) && t.Before(fourPM)
}
