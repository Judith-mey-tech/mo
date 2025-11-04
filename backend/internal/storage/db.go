package storage

import (
    "database/sql"
    "errors"
    "fmt"

    "github.com/mattn/go-sqlite3"
)

type FileRecord struct {
    ID       int64
    Name     string
    Path     string
    SHA256   string
    Size     int64
    Created  string
}

type CheckResult struct {
    FileID    int64   ⁠ json:"file_id" ⁠
    Confidence float64 ⁠ json:"confidence" ⁠
    Reason     string  ⁠ json:"reason" ⁠
    CheckedAt  string  ⁠ json:"checked_at" ⁠
}

func NewSQLite(path string) (*sql.DB, error) {
    db, err := sql.Open("sqlite3", path)
    if err != nil {
        return nil, err
    }
    if err := migrate(db); err != nil {
        db.Close()
        return nil, err
    }
    return db, nil
}

func migrate(db *sql.DB) error {
    stmts := []string{
        `CREATE TABLE IF NOT EXISTS users (
           id INTEGER PRIMARY KEY AUTOINCREMENT,
           username TEXT UNIQUE,
           password TEXT
        );`,
        `CREATE TABLE IF NOT EXISTS files (
           id INTEGER PRIMARY KEY AUTOINCREMENT,
           name TEXT,
           path TEXT,
           sha256 TEXT,
           size INTEGER,
           created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );`,
        `CREATE TABLE IF NOT EXISTS checks (
           id INTEGER PRIMARY KEY AUTOINCREMENT,
           file_id INTEGER,
           confidence REAL,
           reason TEXT,
           checked_at DATETIME DEFAULT CURRENT_TIMESTAMP,
           FOREIGN KEY(file_id) REFERENCES files(id)
        );`,
    }
    for _, s := range stmts {
        if _, err := db.Exec(s); err != nil {
            return err
        }
    }
    return nil
}

// User simple functions (password stored in plaintext here — в demo, никогда так не делайте в реальном проде)

func CreateUser(db *sql.DB, username, password string) error {
    _, err := db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, password)
    if err != nil {
        return err
    }
    return nil
}

func CheckUser(db *sql.DB, username, password string) (bool, error) {
    var pass string
    row := db.QueryRow("SELECT password FROM users WHERE username = ?", username)
    if err := row.Scan(&pass); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return false, nil
        }
        return false, err
    }
    return pass == password, nil
}

// Files

func CreateFileRecord(db *sql.DB, name, path, sha string, size int64) (int64, error) {
    res, err := db.Exec("INSERT INTO files(name, path, sha256, size) VALUES(?, ?, ?, ?)", name, path, sha, size)
    if err != nil {
        return 0, err
    }
    return res.LastInsertId()
}

func GetFileRecord(db *sql.DB, id int64) (*FileRecord, error) {
    row := db.QueryRow("SELECT id, name, path, sha256, size, created_at FROM files WHERE id = ?", id)
    var r FileRecord
    if err := row.Scan(&r.ID, &r.Name, &r.Path, &r.SHA256, &r.Size, &r.Created); err != nil {
        return nil, err
    }
    return &r, nil
}

func SaveCheckResult(db *sql.DB, fileID int64, confidence float64, reason string) error {
    _, err := db.Exec("INSERT INTO checks(file_id, confidence, reason) VALUES(?, ?, ?)", fileID, confidence, reason)
    return err
}

func GetCheckResult(db *sql.DB, fileID int64) (*CheckResult, error) {
    row := db.QueryRow("SELECT file_id, confidence, reason, checked_at FROM checks WHERE file_id = ? ORDER BY checked_at DESC LIMIT 1", fileID)
    var cr CheckResult
    if err := row.Scan(&cr.FileID, &cr.Confidence, &cr.Reason, &cr.CheckedAt); err != nil {
        return nil, err
    }
    return &cr, nil
}
