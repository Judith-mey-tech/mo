package models

import "time"

// Пользователь
type User struct {
    ID       int64  ⁠ json:"id" ⁠
    Username string ⁠ json:"username" ⁠
    Password string ⁠ json:"password,omitempty" ⁠
    Created  time.Time ⁠ json:"created_at" ⁠
}

// Загруженный файл
type File struct {
    ID        int64     ⁠ json:"id" ⁠
    UserID    int64     ⁠ json:"user_id" ⁠
    Name      string    ⁠ json:"name" ⁠
    Path      string    ⁠ json:"path" ⁠
    SHA256    string    ⁠ json:"sha256" ⁠
    Size      int64     ⁠ json:"size" ⁠
    Uploaded  time.Time ⁠ json:"uploaded_at" ⁠
}

// Результат проверки (AI-анализ)
type CheckResult struct {
    ID        int64     ⁠ json:"id" ⁠
    FileID    int64     ⁠ json:"file_id" ⁠
    Result    string    ⁠ json:"result" ⁠
    Confidence float64  ⁠ json:"confidence" ⁠
    CheckedAt time.Time ⁠ json:"checked_at" ⁠
}