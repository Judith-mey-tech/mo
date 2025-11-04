package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "time"

    "authentrack/internal/storage"

    "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("demo-secret-change-me")

type registerReq struct {
    Username string ⁠ json:"username" ⁠
    Password string ⁠ json:"password" ⁠
}

type loginReq struct {
    Username string ⁠ json:"username" ⁠
    Password string ⁠ json:"password" ⁠
}

type tokenResp struct {
    Token string ⁠ json:"token" ⁠
}

func RegisterHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req registerReq
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        if req.Username == "" || req.Password == "" {
            http.Error(w, "username/password required", http.StatusBadRequest)
            return
        }
        if err := storage.CreateUser(db, req.Username, req.Password); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusCreated)
    }
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req loginReq
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        ok, err := storage.CheckUser(db, req.Username, req.Password)
        if err != nil || !ok {
            http.Error(w, "invalid credentials", http.StatusUnauthorized)
            return
        }

        token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
            "sub": req.Username,
            "exp": time.Now().Add(time.Hour * 72).Unix(),
        })
        signed, err := token.SignedString(jwtSecret)
        if err != nil {
            http.Error(w, "token error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(tokenResp{Token: signed})
    }
}