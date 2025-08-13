package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/golang-jwt/jwt/v5"
    "github.com/ZXyoukai/InfraWatch-Backend/internal/http/middleware"
    prismaClient "github.com/ZXyoukai/InfraWatch-Backend/internal/prismaClient"
)

func MeHandler(pc *prismaClient.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        claims := r.Context().Value(middleware.AuthUserKey).(jwt.MapClaims)
        json.NewEncoder(w).Encode(claims)
    }
}