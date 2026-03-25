package handlers

import (
	"auth-project/internal/middleware"
	"auth-project/internal/service"

	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type NoteHandler struct {
	service *service.NoteService
}

func NewNoteHandler(s *service.NoteService) *NoteHandler {
	return &NoteHandler{service: s}
}

func (h *NoteHandler) Create(w http.ResponseWriter, r *http.Request) {

	userID := middleware.GetUserID(r.Context())

	var req struct {
		Text string `json:"text"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	err := h.service.Create(userID, req.Text)
	if err != nil {
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("created"))
}

func (h *NoteHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	userID := middleware.GetUserID(r.Context())

	notes, err := h.service.GetAll(userID)
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(notes)
}

func (h *NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {

	userID := middleware.GetUserID(r.Context())

	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)

	err := h.service.Delete(id, userID)
	if err != nil {
		http.Error(w, "not allowed", http.StatusForbidden)
		return
	}

	w.Write([]byte("deleted"))
}
