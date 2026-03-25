package main

import (
	"log"
	"net/http"

	"auth-project/internal/handlers"
	"auth-project/internal/middleware"
	"auth-project/internal/repository"
	"auth-project/internal/service"

	"github.com/go-chi/chi/v5"
)

func main() {

	connString := "postgres://postgres:password@localhost:5432/auth"

	db, err := repository.NewPostgres(connString)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(db)

	authService := service.NewAuthService(userRepo)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(authService)

	r := chi.NewRouter()

	r.Post("/auth/login", authHandler.Login)
	r.Post("/auth/logout", authHandler.Logout)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)

		r.Get("/me", userHandler.Me)
	})

	log.Println("Server started on :8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}

	noteRepo := repository.NewNoteRepository(db)
	noteService := service.NewNoteService(noteRepo)
	noteHandler := handlers.NewNoteHandler(noteService)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)

		r.Get("/me", userHandler.Me)

		r.Get("/notes", noteHandler.GetAll)
		r.Post("/notes", noteHandler.Create)
		r.Delete("/notes/{id}", noteHandler.Delete)
	})
}
