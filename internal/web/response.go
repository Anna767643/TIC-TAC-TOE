package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// writeJSON - отправляет код ответа и тело ответа пользователю
func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		fmt.Println("Error parse json:", err)
	}
}

// writeError - оправляет код ошибки и информацию об ошибке пользователю
func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(errorDTO{Error: msg})
	if err != nil {
		fmt.Println("Error parse json:", err)
	}
}
