package calculations

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ReverseInput struct {
	Numbers []int `json:"numbers"`
}

func ReverseHandler(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input ReverseInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			logger.Error("Failed to decode JSON", "error", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		logger.Info("Input received", "numbers:", input.Numbers)
		result := Reverse(input.Numbers)
		w.Header().Set("Content-Type", "application/json")
		logger.Info("Reverse calculated", "reversed:", result)
		json.NewEncoder(w).Encode(result)

	}
}

func Reverse(arr []int) []int {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}
