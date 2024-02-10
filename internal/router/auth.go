package router

import (
	"net/http"

	"github.com/rs/zerolog"
	"github.com/tarro-dev/discord-oauth/internal/auth"
)

func (r *Router) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, r.redirectUri, http.StatusFound)
	}
}

func (r *Router) Callback() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Get the code from the query string
		code := req.URL.Query().Get("code")
		if code == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// Request the token
		token, err := r.discord.RequestToken(req.Context(), code)
		if err != nil {
			zerolog.Ctx(req.Context()).Error().Err(err).Msg("failed to request token")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Create a new session
		session := r.sessions.NewSession(req.Context(), token)

		// Set the session cookie
		http.SetCookie(w, session.Cookie())

		// Store the session in the context
		req = req.WithContext(auth.WithContext(req.Context(), session))

		// Redirect to the home page
		http.Redirect(w, req, "/", http.StatusFound)
	}
}

func (r *Router) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Get the session from the context
		session, ok := auth.Ctx(req.Context()).Unwrap()
		if ok {
			// Delete the session from the store
			r.sessions.DeleteSession(req.Context(), session.ID)
		}

		// Delete the session cookie
		http.SetCookie(w, &http.Cookie{
			Name:   "session",
			Value:  "",
			MaxAge: -1,
		})

		// Redirect to the home page
		http.Redirect(w, req, "/", http.StatusFound)
	}
}
