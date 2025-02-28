package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Вызов следующего обработчика
		next.ServeHTTP(w, r)

		// Логирование запроса
		log.Printf(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			time.Since(start),
		)
	})
}

// ValidationError представляет ошибку валидации
type ValidationError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// RespondWithError отправляет ошибку в формате JSON
func RespondWithError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ValidationError{
		Status:  status,
		Message: message,
	})
}
