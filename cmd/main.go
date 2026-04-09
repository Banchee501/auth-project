package main

import (
	"log"
	"net/http"
	"time"

	"auth-project/internal/config"
	"auth-project/internal/handlers"
	"auth-project/internal/middleware"
	"auth-project/internal/repository"
	"auth-project/internal/service"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {

	cfg := config.Load()
	db, err := repository.NewPostgres(cfg.DBUrl)

	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(db)
	refreshRepo := repository.NewRefreshRepository(db)
	noteRepo := repository.NewNoteRepository(db)

	authService := service.NewAuthService(
		userRepo,
		refreshRepo,
		cfg.JWTSecret,
		cfg.AccessTTL,
		cfg.RefreshTTL,
	)
	noteService := service.NewNoteService(noteRepo)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(authService)
	noteHandler := handlers.NewNoteHandler(noteService)

	r := chi.NewRouter()

	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Post("/refresh", authHandler.Refresh)
		r.Post("/logout", authHandler.Logout)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)

		r.Get("/me", userHandler.Me)

		r.Route("/notes", func(r chi.Router) {
			r.Get("/", noteHandler.GetAll)
			r.Post("/", noteHandler.Create)
			r.Delete("/{id}", noteHandler.Delete)
		})
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Println("Server started on :8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
