package chat

import "regexp"

func validateRussianPhoneNumber(number string) bool {
	phoneRegex := regexp.MustCompile(`^\+7\d{10}$`)
	return phoneRegex.MatchString(number)
}
