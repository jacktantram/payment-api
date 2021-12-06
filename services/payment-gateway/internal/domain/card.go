package domain

import (
	"strconv"
	"strings"
	"unicode"
)

// ValidCardNumber Determines whether or not the card number is valid
// as per luhn algorithm
// https://en.wikipedia.org/wiki/Luhn_algorithm
func ValidCardNumber(number string) bool {
	number = strings.ReplaceAll(number, " ", "")
	if len(number) != 16 {
		return false
	}
	digits := make([]int, len(number))
	for i, n := range number {
		if !unicode.IsDigit(n) {
			return false
		}
		number, err := strconv.Atoi(string(n))
		if err != nil {
			return false
		}
		digits[i] = number

	}
	for i := len(digits) - 2; i >= 0; i -= 2 {
		doubled := digits[i] * 2
		if doubled > 9 {
			doubled -= 9
		}
		digits[i] = doubled
	}
	sum := 0
	for _, n := range digits {
		sum += n
	}
	return sum%10 == 0
}
