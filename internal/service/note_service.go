package service

import (
	"auth-project/internal/models"
	"auth-project/internal/repository"
)

type NoteService struct {
	repo *repository.NoteRepository
}

func NewNoteService(r *repository.NoteRepository) *NoteService {
	return &NoteService{repo: r}
}

func (s *NoteService) Create(userID int, text string) error {
	return s.repo.Create(userID, text)
}

func (s *NoteService) GetAll(userID int) ([]models.Note, error) {
	return s.repo.GetAllByUser(userID)
}

func (s *NoteService) Delete(noteID, userID int) error {
	return s.repo.Delete(noteID, userID)
}
