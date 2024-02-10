package config

const (
	discordClientIDKey     = "discord_client_id"
	discordClientSecretKey = "discord_client_secret"
	discordAuthorizeURIKey = "discord_authorize_uri"
	discordCallbackKey     = "discord_callback"
)

// DiscordClientID returns the Discord client ID.
func (c *Config) DiscordClientID() (string, error) {
	if !c.v.IsSet(discordClientIDKey) {
		return "", ErrMissingKey{Key: discordClientIDKey}
	}

	return c.v.GetString(discordClientIDKey), nil
}

// DiscordClientSecret returns the Discord client secret.
func (c *Config) DiscordClientSecret() (string, error) {
	if !c.v.IsSet(discordClientSecretKey) {
		return "", ErrMissingKey{Key: discordClientSecretKey}
	}

	return c.v.GetString(discordClientSecretKey), nil
}

// DiscordAuthorizeURI returns the Discord URI used for OAuth.
func (c *Config) DiscordAuthorizeURI() (string, error) {
	if !c.v.IsSet(discordAuthorizeURIKey) {
		return "", ErrMissingKey{Key: discordAuthorizeURIKey}
	}

	return c.v.GetString(discordAuthorizeURIKey), nil
}

// DiscordCallback returns the Discord callback URI.
func (c *Config) DiscordCallback() (string, error) {
	if !c.v.IsSet(discordCallbackKey) {
		return "", ErrMissingKey{Key: discordCallbackKey}
	}

	return c.v.GetString(discordCallbackKey), nil
}
