package httpjson

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func Write(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("encode JSON: %v", err)
	}
}
func Error(w http.ResponseWriter, status int, message string) {
	Write(w, status, map[string]string{"error": message})
}
func Decode(w http.ResponseWriter, r *http.Request, v any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 32*1024)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		Error(w, 400, "Invalid JSON request")
		return false
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		Error(w, 400, "Only one JSON object is allowed")
		return false
	}
	return true
}
