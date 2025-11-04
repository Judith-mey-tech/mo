package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"

    "authentrack/internal/storage"
)

type UserProfile struct {
    Username string ⁠ json:"username" ⁠
    Files    []storage.FileRecord ⁠ json:"files" ⁠
}

// GET /user/{username}
func GetUserHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        username := r.URL.Query().Get("username")
        if username == "" {
            http.Error(w, "username required", http.StatusBadRequest)
            return
        }

        files, err := storage.GetUserFiles(db, username)
        if err != nil {
            http.Error(w, "error fetching files", http.StatusInternalServerError)
            return
        }

        profile := UserProfile{
            Username: username,
            Files:    files,
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(profile)
    }
}