package handlers

import (
	"auth-project/internal/service"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	userAgent := r.UserAgent()
	ip := r.RemoteAddr

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	access, refresh, err := h.service.Login(
		req.Email,
		req.Password,
		userAgent,
		ip,
	)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    access,
		HttpOnly: true,
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		HttpOnly: true,
		Path:     "/",
	})

	w.Write([]byte("ok"))
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("refresh_token")
	if err == nil {
		_ = h.service.DeleteRefreshToken(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "access_token",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:   "refresh_token",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	})

	w.Write([]byte("logged out"))
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "no refresh token", http.StatusUnauthorized)
		return
	}

	access, refresh, err := h.service.Refresh(cookie.Value)
	if err != nil {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    access,
		HttpOnly: true,
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		HttpOnly: true,
		Path:     "/",
	})

	w.Write([]byte("refreshed"))
}
