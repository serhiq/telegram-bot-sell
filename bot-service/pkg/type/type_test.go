package _type

import (
	"testing"
)

func TestFormatPrice(t *testing.T) {
	testCases := []struct {
		priceInKopeks uint64
		expected      string
	}{
		{100, "1 ₽"},
		{0, "0 ₽"},
		{150, "1.50 ₽"},
		{1011, "10.11 ₽"},
		{999999, "9999.99 ₽"},
	}

	for _, testCase := range testCases {
		result := FormatPrice(testCase.priceInKopeks)
		if result != testCase.expected {
			t.Errorf("Unexpected result: expected %s, got %s", testCase.expected, result)
		}
	}
}
