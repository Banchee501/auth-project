package handlers

import (
	"auth-project/internal/middleware"
	"auth-project/internal/service"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	service *service.AuthService
}

func NewUserHandler(s *service.AuthService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {

	userID := middleware.GetUserID(r.Context())

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	resp := map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
	}

	json.NewEncoder(w).Encode(resp)
}
