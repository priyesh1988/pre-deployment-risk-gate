package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/yourname/guardrail-saas/internal/storage"
)

func OrgRiskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		org := parts[len(parts)-2]
		avg, _ := storage.OrgRisk(db, org)
		json.NewEncoder(w).Encode(map[string]float64{"average_risk": avg})
	}
}
