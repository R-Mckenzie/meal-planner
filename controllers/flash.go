package controllers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func setAlertData(w http.ResponseWriter, message string, status string) {
	setFlash(w, "alert_message", []byte(message))
	setFlash(w, "alert_type", []byte(status))
}

func getAlertData(w http.ResponseWriter, r *http.Request) (string, string, error) {
	message, err := getFlash(w, r, "alert_message")
	status, err := getFlash(w, r, "alert_type")
	if err != nil {
		return "", "", err
	}
	return string(message), string(status), nil
}

func setFlash(w http.ResponseWriter, name string, value []byte) {
	c := &http.Cookie{Name: name, Value: encode(value)}
	http.SetCookie(w, c)
}

func getFlash(w http.ResponseWriter, r *http.Request, name string) ([]byte, error) {
	c, err := r.Cookie(name)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return nil, nil
		}
		return nil, fmt.Errorf("in getFlash: %w", err)
	}

	value, err := decode(c.Value)
	if err != nil {
		return nil, err
	}

	dc := &http.Cookie{Name: name, MaxAge: -1, Expires: time.Unix(1, 0)}
	http.SetCookie(w, dc)
	return value, nil
}

func encode(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

func decode(src string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(src)
}
