package auth

import (
	"net/http"
)

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "no refresh token", http.StatusUnauthorized)
		return
	}

	access, refresh, err := h.service.Refresh(cookie.Value)
	if err != nil {
		http.Error(w, "invalid", http.StatusUnauthorized)
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
