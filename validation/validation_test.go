package validation

import (
	"testing"

	"github.com/R-Mckenzie/meal-planner/assert"
)

func TestLongEnough(t *testing.T) {
	tests := []struct {
		title    string
		expected bool
		s        string
		l        int
	}{
		{title: "0 length and 0 minimum", expected: true, s: "", l: 0},
		{title: "0 length and 1 minimum", expected: false, s: "", l: 1},
		{title: "1 length and 0 minimum", expected: true, s: "a", l: 0},
		{title: "1 length and 1 minimum", expected: true, s: "a", l: 1},
		{title: "5 length and 1 minimum", expected: true, s: "aaaaa", l: 1},
		{title: "5 length and 5 minimum", expected: true, s: "aaaaa", l: 5},
		{title: "5 length and 8 minimum", expected: false, s: "aaaaa", l: 8},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			assert.Equals(t, LongEnough(test.s, test.l), test.expected)
		})

	}
}

func TestIsEmail(t *testing.T) {
	tests := []struct {
		expected bool
		email    string
	}{
		{expected: false, email: "abcde"},
		{expected: false, email: "test@"},
		{expected: true, email: "email@email"},
		{expected: true, email: "email@email.com"},
		{expected: true, email: "a@b.c"},
		{expected: false, email: ""},
	}

	for _, test := range tests {
		t.Run(test.email, func(t *testing.T) {
			assert.Equals(t, IsEmail(test.email), test.expected)
		})
	}
}

func TestPasswordCheck(t *testing.T) {
	t.Run("Password is < 8 characters", func(t *testing.T) {
		v, f := PasswordCheck("1234")
		assert.Equals(t, v, false)
		assert.Equals(t, f[0], "Passwords must be at least 8 characters in length.")
	})

	t.Run("Password is = 8 characters", func(t *testing.T) {
		v, f := PasswordCheck("12345678")
		assert.Equals(t, v, true)
		assert.Equals(t, len(f), 0)
	})
}
