package main

import (
    "context"
    "log"
    "net/http"
    "os"

    "github.com/go-chi/chi/v5"
    "github.com/joho/godotenv"
    "github.com/youruser/restful-go-prisma-jwt/internal/http/handlers"
    "github.com/youruser/restful-go-prisma-jwt/internal/http/middleware"
    "github.com/youruser/restful-go-prisma-jwt/internal/prismaClient"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Println(".env not loaded, relying on env vars")
    }

    ctx := context.Background()

    client, err := prismaClient.NewClient(ctx)
    if err != nil {
        log.Fatalln("prisma client init:", err)
    }
    defer client.Disconnect()

    r := chi.NewRouter()

    r.Post("/api/auth/register", handlers.RegisterHandler(client))
    r.Post("/api/auth/login", handlers.LoginHandler(client))

    // protected routes
    r.Group(func(r chi.Router) {
        r.Use(middleware.JWTAuth)
        r.Get("/api/me", handlers.MeHandler(client))
        r.Route("/api/posts", func(r chi.Router) {
            r.Get("/", handlers.ListPosts(client))
            r.Post("/", handlers.CreatePost(client))
            r.Get("/{id}", handlers.GetPost(client))
            r.Put("/{id}", handlers.UpdatePost(client))
            r.Delete("/{id}", handlers.DeletePost(client))
        })
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Println("listening on :" + port)
    http.ListenAndServe(":"+port, r)
}