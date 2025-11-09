package utils

import (
	"regexp"
	"unicode"
)

func CheckPasswordVulnerabilities(password string) string {
	if len(password) < 8 {
		return "Password must be at least 8 characters long"
	}

	hasLetter := false
	hasDigit := false
	for _, char := range password {
		switch {
		case unicode.IsLetter(char):
			hasLetter = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	if !hasLetter {
		return "Password must contain letters"
	}

	if !hasDigit {
		return "Password must contain digits"
	}

	return ""
}

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
