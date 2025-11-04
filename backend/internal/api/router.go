package api

import (
	"database/sql"
	"net/http"

	"authentrack/internal/api/handlers"

	"github.com/go-chi/chi/v5"
)

func NewRouter(db *sql.DB) http.Handler {
    r := chi.NewRouter()

    // public
    r.Post("/auth/register", handlers.RegisterHandler(db))
    r.Post("/auth/login", handlers.LoginHandler(db))

    r.Post("/upload", handlers.UploadHandler(db)) // для простоты без auth
    r.Get("/file/{id}", handlers.GetFileHandler(db))
    r.Get("/check/{id}", handlers.CheckResultHandler(db))

    // health
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("ok"))
    })

    return r
}
