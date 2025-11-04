package storage

import (
    "crypto/sha256"
    "encoding/hex"
    "io"
    "os"
    "path/filepath"
)

// SaveUploadedFile сохраняет файл и возвращает путь и хэш
func SaveUploadedFile(file io.Reader, filename, uploadDir string) (string, string, int64, error) {
    os.MkdirAll(uploadDir, os.ModePerm)

    // Сохраняем файл
    dstPath := filepath.Join(uploadDir, filename)
    dst, err := os.Create(dstPath)
    if err != nil {
        return "", "", 0, err
    }
    defer dst.Close()

    hash := sha256.New()
    size, err := io.Copy(io.MultiWriter(dst, hash), file)
    if err != nil {
        return "", "", 0, err
    }

    sha := hex.EncodeToString(hash.Sum(nil))
    return dstPath, sha, size, nil
}