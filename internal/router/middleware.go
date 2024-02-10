package router

import (
	"net/http"

	"github.com/rs/zerolog"
	"github.com/tarro-dev/discord-oauth/internal/auth"
)

// LoggerMiddleware returns a new logger middleware.
// It adds a logger to the context.
func LoggerMiddleware(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = logger.WithContext(ctx)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// SessionMiddleware returns a new session middleware.
func (r *Router) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Get the session cookie
		cookie, err := req.Cookie("session")
		if err != nil {
			// No session cookie
			next.ServeHTTP(w, req)
			return
		}

		// Get the session from the store
		session, ok := r.sessions.GetSession(req.Context(), cookie.Value).Unwrap()
		if !ok {
			// Invalid session
			http.SetCookie(w, &http.Cookie{
				Name:   "session",
				Value:  "",
				MaxAge: -1,
			})
			next.ServeHTTP(w, req)
			return
		}

		// Store the session in the context
		req = req.WithContext(auth.WithContext(req.Context(), &session))

		zerolog.Ctx(req.Context()).Info().Str("session", session.ID).Msg("session found")

		next.ServeHTTP(w, req)
	})
}
