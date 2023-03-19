package chat

import (
	"testing"
)

func TestValidateRussianPhoneNumber(t *testing.T) {
	validNumbers := []string{"+79111234567", "+71234567890"}
	invalidNumbers := []string{"1234567890", "+7911123456", "+712345678901"}

	for _, number := range validNumbers {
		if !validateRussianPhoneNumber(number) {
			t.Errorf("Error validating valid phone number %s", number)
		}
	}

	for _, number := range invalidNumbers {
		if validateRussianPhoneNumber(number) {
			t.Errorf("Error validating invalid phone number %s", number)
		}
	}
}
