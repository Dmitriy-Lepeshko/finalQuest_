package api

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("my-secret-key-for-todo-app")

type TokenClaims struct {
	PasswordHash string `json:"password_hash"`
	jwt.RegisteredClaims
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJson(w, map[string]string{"error": "invalid JSON"})
		return
	}

	expectedPass := os.Getenv("TODO_PASSWORD")
	if expectedPass == "" {
		writeJson(w, map[string]string{"error": "authentication disabled"})
		return
	}

	if subtle.ConstantTimeCompare([]byte(req.Password), []byte(expectedPass)) != 1 {
		writeJson(w, map[string]string{"error": "Неверный пароль"})
		return
	}

	expirationTime := time.Now().Add(8 * time.Hour)
	claims := &TokenClaims{
		PasswordHash: generateHash(expectedPass),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		writeJson(w, map[string]string{"error": "internal server error"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		MaxAge:   int(8 * time.Hour.Seconds()),
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
	})

	writeJson(w, map[string]string{"token": tokenString})
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if pass == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		tknStr := cookie.Value
		claims := &TokenClaims{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !tkn.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if subtle.ConstantTimeCompare([]byte(claims.PasswordHash), []byte(generateHash(pass))) != 1 {
			http.Error(w, "Token expired or invalid", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func generateHash(password string) string {
	return password
}
