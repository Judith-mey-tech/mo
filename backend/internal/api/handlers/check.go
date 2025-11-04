package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "strconv"

    "authentrack/internal/storage"
    "github.com/go-chi/chi/v5"
)

// GET /check/{id}
func CheckResultHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        idStr := chi.URLParam(r, "id")
        id, _ := strconv.ParseInt(idStr, 10, 64)
        res, err := storage.GetCheckResult(db, id)
        if err != nil {
            http.Error(w, "not found", http.StatusNotFound)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(res)
    }
}