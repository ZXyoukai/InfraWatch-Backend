package handlers

import (
    "context"
    "encoding/json"
    "net/http"
    "os"
    "time"

    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v5"
    prismaClient "github.com/youruser/restful-go-prisma-jwt/internal/prismaClient"
    prisma "github.com/youruser/restful-go-prisma-jwt/prisma/client"
)

func RegisterHandler(pc *prismaClient.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var in struct{ Email, Name, Password string }
        if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
            http.Error(w, "bad request", http.StatusBadRequest)
            return
        }
        hash, _ := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
        user, err := pc.Client.User.CreateOne(
            prisma.User.Email.Set(in.Email),
            prisma.User.Name.Set(&in.Name),
            prisma.User.Password.Set(string(hash)),
        ).Exec(context.Background())
        if err != nil {
            http.Error(w, "unable to create user", http.StatusInternalServerError)
            return
        }
        user.Password = ""
        json.NewEncoder(w).Encode(user)
    }
}

func LoginHandler(pc *prismaClient.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var in struct{ Email, Password string }
        if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
            http.Error(w, "bad request", http.StatusBadRequest)
            return
        }
        user, err := pc.Client.User.FindUnique(
            prisma.User.Email.Equals(in.Email),
        ).Exec(context.Background())
        if err != nil || user == nil {
            http.Error(w, "invalid credentials", http.StatusUnauthorized)
            return
        }
        if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password)) != nil {
            http.Error(w, "invalid credentials", http.StatusUnauthorized)
            return
        }
        // create JWT
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
            "sub": user.ID,
            "email": user.Email,
            "role": user.Role,
            "exp": time.Now().Add(24 * time.Hour).Unix(),
        })
        secret := os.Getenv("JWT_SECRET")
        ss, _ := token.SignedString([]byte(secret))
        json.NewEncoder(w).Encode(map[string]string{"token": ss})
    }
}