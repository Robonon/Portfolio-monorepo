package calculations

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type SumInput struct {
	Numbers []int `json:"numbers"`
}

func SumHandler(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input SumInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			logger.Error("Failed to decode JSON", "error", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		logger.Info("Input received", "numbers:", input.Numbers)
		sum := Sum(input.Numbers)
		logger.Info("Sum calculated", "sum:", sum)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"sum": sum})
	}
}

// 2. Calculate the sum of all elements
func Sum(arr []int) int {
	sum := 0
	for _, v := range arr {
		sum += v
	}
	return sum
}
