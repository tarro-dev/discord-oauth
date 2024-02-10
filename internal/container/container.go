package container

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/tarro-dev/discord-oauth/internal/auth"
	"github.com/tarro-dev/discord-oauth/internal/config"
	"github.com/tarro-dev/discord-oauth/internal/router"
	"github.com/tarro-dev/discord-oauth/internal/user"
	"github.com/tarro-dev/discord-oauth/pkg/singleton"
)

// Container is the container that holds the services.
type Container struct {
	config *config.Config

	router  singleton.Singleton[*router.Router]
	session singleton.Singleton[*auth.SessionStore]
	discord singleton.Singleton[*auth.Discord]
}

// New returns a new container.
func New(config *config.Config) *Container {
	return &Container{
		config:  config,
		router:  singleton.New[*router.Router](),
		session: singleton.New[*auth.SessionStore](),
		discord: singleton.New[*auth.Discord](),
	}
}

// Router returns the http router.
func (c *Container) Router(ctx context.Context) (*router.Router, error) {
	return c.router.GetOr(func() (*router.Router, error) {
		auth, err := c.Discord(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get discord service: %w", err)
		}

		sessions, err := c.SessionStore(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get session store: %w", err)
		}

		uri, err := c.config.DiscordAuthorizeURI()
		if err != nil {
			return nil, fmt.Errorf("failed to get discord uri: %w", err)
		}

		zerolog.Ctx(ctx).Info().Str("authorize_uri", uri).Msg("creating router")

		return router.New(router.RouterParams{
			Discord:     auth,
			Sessions:    sessions,
			Users:       c.UserStore(),
			RedirectUri: uri,
		}), nil
	})
}

// SessionStore returns the session store.
func (c *Container) SessionStore(ctx context.Context) (*auth.SessionStore, error) {
	return c.session.GetOr(func() (*auth.SessionStore, error) {
		discord, err := c.Discord(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to get discord service: %w", err)
		}

		zerolog.Ctx(ctx).Info().Msg("creating session store")

		return auth.NewSessionStore(discord), nil
	})
}

// Discord returns the discord auth service.
func (c *Container) Discord(ctx context.Context) (*auth.Discord, error) {
	return c.discord.GetOr(func() (*auth.Discord, error) {
		clientID, err := c.config.DiscordClientID()
		if err != nil {
			return nil, fmt.Errorf("failed to get discord client id: %w", err)
		}

		clientSecret, err := c.config.DiscordClientSecret()
		if err != nil {
			return nil, fmt.Errorf("failed to get discord client secret: %w", err)
		}

		uri, err := c.config.DiscordCallback()
		if err != nil {
			return nil, fmt.Errorf("failed to get discord uri: %w", err)
		}

		zerolog.Ctx(ctx).Info().Str("client_id", clientID).Str("callback", uri).Msg("creating discord service")

		return auth.NewDiscord(clientID, clientSecret, uri), nil
	})
}

// UserStore returns the user store.
func (c *Container) UserStore() *user.UserStore {
	return user.NewUserStore()
}
