package user

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/tarro-dev/discord-oauth/internal/auth"
)

// User is the user information.
type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// UserStore is the user store.
type UserStore struct{}

// NewUserStore returns a new user store.
func NewUserStore() *UserStore {
	return &UserStore{}
}

// GetUser returns the user based on the session.
func (s *UserStore) GetUser(session auth.Session) (*User, error) {
	ds, err := discordgo.New("Bearer " + session.Token.Access)
	if err != nil {
		return nil, fmt.Errorf("failed to create discord session: %w", err)
	}

	u, err := ds.User("@me")
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &User{
		ID:     u.ID,
		Name:   u.Username,
		Avatar: u.AvatarURL("512"),
	}, nil
}
