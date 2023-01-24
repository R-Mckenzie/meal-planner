package validation

import (
	"net/mail"
	"unicode"
)

// Returns true if s is >= min, false otherwise
func LongEnough(s string, min int) bool {
	return len(s) >= min
}

// Returns true if s contains any unicode symbols, false otherwise
func ContainsSymbol(s string) bool {
	for _, c := range s {
		if unicode.IsSymbol(c) {
			return true
		}
	}
	return false
}

// Returns true if s contains any numbers, false otherwise
func ContainsNumber(s string) bool {
	for _, c := range s {
		if unicode.IsNumber(c) {
			return true
		}
	}
	return false
}

// Returns true if s is a valid RFC5322 spec email address
func IsEmail(s string) bool {
	_, err := mail.ParseAddress(s)
	return err == nil
}

// If the password contains symbols, numbers, and is over 8 chars returns true, nil
// If any checks fail, returns false, with a string slice of all faults
func PasswordCheck(p string) (bool, []string) {
	e := make([]string, 0)
	valid := true
	if !ContainsSymbol(p) || !ContainsNumber(p) {
		e = append(e, "Passwords must contain both symbols and numbers.")
		valid = false
	}
	if !LongEnough(p, 8) {
		e = append(e, "Passwords must be at least 8 characters in length.")
		valid = false
	}

	if valid {
		e = nil
	}
	return valid, e
}
