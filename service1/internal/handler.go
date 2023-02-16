package internal

import (
	"encoding/json"
	"math/rand"
	"net/http"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	saltLen = 12
)

type Reply struct {
	Salt string `json:"salt"`
}

func salt() string {
	b := make([]byte, 0, saltLen)
	for len(b) < saltLen {
		b = append(b, charset[rand.Intn(len(charset))])
	}
	return string(b)
}

func GenerateSalt(w http.ResponseWriter, r *http.Request) {
	var resp Reply
	resp.Salt = salt()

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
