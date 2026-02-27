package handlers


import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
)

type Request struct {
	Data string `json:"data"`
}

type Response struct {
	Word []string `json:"word"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func SortHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		writeError(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req Request
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		writeError(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.Data == "" {
		writeError(w, `Missing or empty "data" field`, http.StatusBadRequest)
		return
	}

	if len(req.Data) > 1000 {
		writeError(w, `"data" field exceeds maximum length of 1000 characters`, http.StatusBadRequest)
		return
	}

	chars := strings.Split(req.Data, "")

	sort.Slice(chars, func(i, j int) bool {
		return strings.ToLower(chars[i]) < strings.ToLower(chars[j])
	})

	resp := Response{Word: chars}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func writeError(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: msg})
}