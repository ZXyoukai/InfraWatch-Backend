package middleware

import (
    "context"
    "net/http"
    "os"
    "strings"

    "github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const AuthUserKey ctxKey = "authUser"

func JWTAuth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        auth := r.Header.Get("Authorization")
        if auth == "" {
            http.Error(w, "missing auth", http.StatusUnauthorized)
            return
        }
        parts := strings.SplitN(auth, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, "invalid auth header", http.StatusUnauthorized)
            return
        }
        tokenStr := parts[1]
        secret := os.Getenv("JWT_SECRET")
        token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
            return []byte(secret), nil
        })
        if err != nil || !token.Valid {
            http.Error(w, "invalid token", http.StatusUnauthorized)
            return
        }
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            http.Error(w, "invalid claims", http.StatusUnauthorized)
            return
        }
        // attach user data to context
        ctx := context.WithValue(r.Context(), AuthUserKey, claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}