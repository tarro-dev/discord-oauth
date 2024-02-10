package router

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/tarro-dev/discord-oauth/internal/auth"
	"github.com/tarro-dev/discord-oauth/internal/user"
	"github.com/tarro-dev/discord-oauth/pkg/option"
	"github.com/tarro-dev/discord-oauth/templates"
)

// Router is the http router.
type Router struct {
	discord     *auth.Discord
	sessions    *auth.SessionStore
	users       *user.UserStore
	redirectUri string
}

// RouterParams is the params to create a new router.
type RouterParams struct {
	Discord     *auth.Discord
	Sessions    *auth.SessionStore
	Users       *user.UserStore
	RedirectUri string
}

// New returns a new router.
func New(params RouterParams) *Router {
	return &Router{
		discord:     params.Discord,
		sessions:    params.Sessions,
		users:       params.Users,
		redirectUri: params.RedirectUri,
	}
}

// Handler returns a new handler.
func (router *Router) Handler(ctx context.Context) http.Handler {
	r := chi.NewRouter()
	r.Use(LoggerMiddleware(*zerolog.Ctx(ctx)))
	r.Use(router.SessionMiddleware)

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		user := router.RetrieveUser(r.Context())

		if err := templates.Index(user.UnwrapPtr()).Render(r.Context(), w); err != nil {
			zerolog.Ctx(r.Context()).Error().Err(err).Msg("failed to render index")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	})
	r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
		session, ok := router.LoginPage(w, r)
		if !ok {
			return
		}

		user, err := router.users.GetUser(*session)
		if err != nil {
			zerolog.Ctx(r.Context()).Error().Err(err).Msg("failed to get user")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err := templates.ProfilePage(user).Render(r.Context(), w); err != nil {
			zerolog.Ctx(r.Context()).Error().Err(err).Msg("failed to render profile page")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	})

	// Auth
	r.Get("/login", router.Login())
	r.Get("/logout", router.Logout())
	r.Get("/callback", router.Callback())

	return r
}

// LoginPage returns the login page if the user is not authenticated.
func (router *Router) LoginPage(w http.ResponseWriter, r *http.Request) (*auth.Session, bool) {
	session, ok := auth.Ctx(r.Context()).Unwrap()
	if !ok {
		if err := templates.Login().Render(r.Context(), w); err != nil {
			zerolog.Ctx(r.Context()).Error().Err(err).Msg("failed to render login page")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return &auth.Session{}, false
	}

	return session, true
}

// RetrieveUser returns the user from the session as an option.
func (router *Router) RetrieveUser(ctx context.Context) option.Option[user.User] {
	session, ok := auth.Ctx(ctx).Unwrap()
	if !ok {
		return option.None[user.User]()
	}

	u, err := router.users.GetUser(*session)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to get user")
		return option.None[user.User]()
	}

	return option.Some(*u)
}
