package validation

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// A key value store of errors
type Errors map[string][]string

// Any return true if there is any error.
func (e Errors) Any() bool {
	return len(e) > 0
}

// Add adds an error for a specific field
func (e Errors) Add(field string, msg string) {
	if _, ok := e[field]; !ok {
		e[field] = []string{}
	}
	e[field] = append(e[field], msg)
}

// Get returns all the errors for the given field.
func (e Errors) Get(field string) []string {
	return e[field]
}

// Has returns true whether the given field has any errors.
func (e Errors) Has(field string) bool {
	return len(e[field]) > 0
}

// Regular expression defining a valid email
var emailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// Returns true if the value is at least the minumum length
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// Returns true if the value is a valid email
func IsEmail(value string) bool {
	return emailRX.MatchString(value)
}

func Validate(field, errMsg string, validationFunc bool, errors Errors) {
	if !validationFunc {
		errors.Add(field, errMsg)
	}
}
