package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/R-Mckenzie/meal-planner/assert"
)

func TestSecureHeaders(t *testing.T) {
	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal()
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	secureHeaders(next).ServeHTTP(rr, r)
	res := rr.Result()

	//assert.Equals(t, res.Header.Get("Content-Security-Policy"), "default-src 'self'; script-src 'self' 'unsafe-eval'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
	assert.Equals(t, res.Header.Get("Referrer-Policy"), "origin-when-cross-origin")
	assert.Equals(t, res.Header.Get("X-Content-Type-Options"), "nosniff")
	assert.Equals(t, res.Header.Get("X-Frame-Options"), "deny")
	assert.Equals(t, res.Header.Get("X-XSS-Protection"), "0")
	assert.Equals(t, res.StatusCode, http.StatusOK)

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	assert.Equals(t, string(body), "OK")
}
