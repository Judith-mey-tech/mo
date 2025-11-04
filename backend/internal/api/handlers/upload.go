package handlers

import (
    "crypto/sha256"
    "database/sql"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "time"

    "authentrack/internal/core"
    "authentrack/internal/storage"
)

type uploadResp struct {
    ID      int64  `json:"id"`
    SHA     string `json:"sha256"`
    Message string `json:"message"`
}

func UploadHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Максимум 200MB для демонстрации
        r.Body = http.MaxBytesReader(w, r.Body, 200<<20)
        if err := r.ParseMultipartForm(200 << 20); err != nil {
            http.Error(w, "file too big or parse error", http.StatusBadRequest)
            return
        }

        file, header, err := r.FormFile("file")
        if err != nil {
            http.Error(w, "file required", http.StatusBadRequest)
            return
        }
        defer file.Close()

        // Сохранить во временный файл и посчитать SHA256
        filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(header.Filename))
        dstPath := filepath.Join("uploads", filename)
        dst, err := os.Create(dstPath)
        if err != nil {
            http.Error(w, "cannot save file", http.StatusInternalServerError)
            return
        }
        defer dst.Close()

        hasher := sha256.New()
        mw := io.MultiWriter(dst, hasher)
        if _, err := io.Copy(mw, file); err != nil {
            http.Error(w, "save error", http.StatusInternalServerError)
            return
        }
        sha := hex.EncodeToString(hasher.Sum(nil))

        // Сохранить запись в БД
        id, err := storage.CreateFileRecord(db, header.Filename, dstPath, sha, header.Size)
        if err != nil {
            http.Error(w, "db save error", http.StatusInternalServerError)
            return
        }

        // Выполнить простую проверку (core)
        result := core.SimpleVerify(dstPath, header.Filename, sha)

        // Сохранить результаты проверки
        if err := storage.SaveCheckResult(db, id, result.Confidence, result.Reason); err != nil {
            // не фатально
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(uploadResp{
            ID: id,
            SHA: sha,
            Message: fmt.Sprintf("uploaded, confidence: %.2f", result.Confidence),
        })
    }
}

func GetFileHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        idStr := chi.URLParam(r, "id")
        id, _ := strconv.ParseInt(idStr, 10, 64)
        rec, err := storage.GetFileRecord(db, id)
        if err != nil {
            http.Error(w, "not found", http.StatusNotFound)
            return
        }
        http.ServeFile(w, r, rec.Path)
    }
}

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
