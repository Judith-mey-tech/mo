package main

import (
    "log"
    "net/http"
    "os"

    "authentrack/internal/api"
    "authentrack/internal/storage"
)

func main() {
    // Создать папку uploads
    if err := os.MkdirAll("uploads", 0755); err != nil {
        log.Fatalf("failed to create uploads dir: %v", err)
    }

    // Инициализировать БД (создаст файл authentrack.db)
    db, err := storage.NewSQLiteDB("authentrack.db")
    if err != nil {
        log.Fatalf("db init: %v", err)
    }
    defer db.Close()

    r := api.NewRouter(db)

    srv := &http.Server{
        Addr:    ":8080",
        Handler: r,
    }

    log.Println("Server running on :8080")
    if err := srv.ListenAndServe(); err != nil {
        log.Fatalf("server error: %v", err)
    }
}