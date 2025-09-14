package calculations

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type MaxInput struct {
	Numbers []int `json:"numbers"`
}

func MaxHandler(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input MaxInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			logger.Error("Failed to decode JSON", "error", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		logger.Info("Input received", "numbers:", input.Numbers)
		result := Max(input.Numbers)
		w.Header().Set("Content-Type", "application/json")
		logger.Info("Max calculated", "max:", result)
		json.NewEncoder(w).Encode(map[string]int{"max": result})
	}
}

func Max(arr []int) int {
	if len(arr) == 0 {
		panic("empty array")
	}
	max := arr[0]
	for _, v := range arr {
		if v > max {
			max = v
		}
	}
	return max
}
