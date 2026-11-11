package health

import (
	"context"
	"eneba-marketplace/backend/internal/httpjson"
	"net/http"
	"time"
)

type Pinger interface{ PingContext(context.Context) error }

type Response struct {
	Status   string `json:"status"`
	Service  string `json:"service"`
	Database string `json:"database"`
}

func Handler(db Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()
		result := Response{Status: "ok", Service: "eneba-api", Database: "connected"}
		if err := db.PingContext(ctx); err != nil {
			result.Status = "error"
			result.Database = "unavailable"
			httpjson.Write(w, http.StatusServiceUnavailable, result)
			return
		}
		httpjson.Write(w, http.StatusOK, result)
	}
}
