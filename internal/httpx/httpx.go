package httpx

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v) 
}

func JSONError(w http.ResponseWriter, status int, msg string) {
	type Err struct{ Error string `json:"error"` }
	JSON(w, status, Err{Error: msg})
}