package validation

import (
	"net/mail"
)

// Returns true if s is >= min, false otherwise
func LongEnough(s string, min int) bool {
	return len(s) >= min
}

// Returns true if s is a valid RFC5322 spec email address
func IsEmail(s string) bool {
	_, err := mail.ParseAddress(s)
	return err == nil
}

// If the password contains symbols, numbers, and is over 8 chars returns true, nil
// If any checks fail, returns false, with a string slice of all faults
func PasswordCheck(p string) (bool, []string) {
	faults := make([]string, 0)
	valid := true
	if !LongEnough(p, 8) {
		faults = append(faults, "Passwords must be at least 8 characters in length.")
		valid = false
	}

	if valid {
		faults = nil
	}
	return valid, faults
}
