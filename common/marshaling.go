package common

import (
	"encoding/json"
	"net/http"
)

func Encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		return err
	}
	return nil
}

func Decode[T any](r *http.Request) (T, error) {
	var v T
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		return v, err
	}
	return v, nil
}
