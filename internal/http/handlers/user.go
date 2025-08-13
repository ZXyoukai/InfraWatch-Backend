package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/golang-jwt/jwt/v5"
    "github.com/youruser/restful-go-prisma-jwt/internal/http/middleware"
    prismaClient "github.com/youruser/restful-go-prisma-jwt/internal/prismaClient"
)

func MeHandler(pc *prismaClient.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        claims := r.Context().Value(middleware.AuthUserKey).(jwt.MapClaims)
        json.NewEncoder(w).Encode(claims)
    }
}