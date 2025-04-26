package report

import (
	"encoding/json"
	"math/rand"
	"net/http"
)

type Report struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func ReportsHandler(w http.ResponseWriter, r *http.Request) {
	reports := []Report{
		{ID: rand.Intn(1000), Name: "Report 1"},
		{ID: rand.Intn(1000), Name: "Report 2"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"reports": reports,
	})
}
