package handlers

import (
    "context"
    "encoding/json"
    "net/http"

    "github.com/go-chi/chi/v5"
    prismaClient "github.com/ZXyoukai/InfraWatch-Backend/internal/prismaClient"
    prisma "github.com/ZXyoukai/InfraWatch-Backend/prisma/client"
)

func ListPosts(pc *prismaClient.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        posts, _ := pc.Client.Post.FindMany().Exec(context.Background())
        json.NewEncoder(w).Encode(posts)
    }
}

func CreatePost(pc *prismaClient.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var in struct{ Title, Body string }
        json.NewDecoder(r.Body).Decode(&in)
        p, _ := pc.Client.Post.CreateOne(
            prisma.Post.Title.Set(in.Title),
            prisma.Post.Body.Set(in.Body),
        ).Exec(context.Background())
        json.NewEncoder(w).Encode(p)
    }
}

func GetPost(pc *prismaClient.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id := chi.URLParam(r, "id")
        p, _ := pc.Client.Post.FindUnique(
            prisma.Post.ID.Equals(id),
        ).Exec(context.Background())
        if p == nil { http.NotFound(w, r); return }
        json.NewEncoder(w).Encode(p)
    }
}

func UpdatePost(pc *prismaClient.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id := chi.URLParam(r, "id")
        var in struct{ Title, Body string; Published *bool }
        json.NewDecoder(r.Body).Decode(&in)
        up := pc.Client.Post.FindUnique(prisma.Post.ID.Equals(id)).Update(
            prisma.Post.Title.Set(in.Title),
            prisma.Post.Body.Set(in.Body),
        ).Exec(context.Background())
        json.NewEncoder(w).Encode(up)
    }
}

func DeletePost(pc *prismaClient.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id := chi.URLParam(r, "id")
        pc.Client.Post.FindUnique(prisma.Post.ID.Equals(id)).Delete().Exec(context.Background())
        w.WriteHeader(http.StatusNoContent)
    }
}
