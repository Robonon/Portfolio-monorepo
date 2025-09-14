package calculations

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type CountUniqueInput struct {
	Numbers []int `json:"numbers"`
}

func CountUniqueHandler(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input CountUniqueInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			logger.Error("Failed to decode JSON", "error", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		logger.Info("Input received", "numbers:", input.Numbers)
		result := CountUnique(input.Numbers)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"unique_count": result})
	}
}

// 5. Count the number of unique elements
func CountUnique(arr []int) int {
	seen := make(map[int]struct{})
	for _, v := range arr {
		seen[v] = struct{}{}
	}
	return len(seen)
}
