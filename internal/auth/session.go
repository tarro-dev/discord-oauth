package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/tarro-dev/discord-oauth/pkg/option"
)

const expirationLimit time.Duration = 30 * time.Minute

type SessionStore struct {
	discord *Discord
	store   *store
}

func NewSessionStore(discord *Discord) *SessionStore {
	return &SessionStore{
		store:   newStore(),
		discord: discord,
	}
}

func (s *SessionStore) NewSession(ctx context.Context, token Token) *Session {
	id := uuid.New().String()

	session := &Session{
		ID:    id,
		Token: token,
	}

	s.store.store(*session)

	return session
}

func (s *SessionStore) GetSession(ctx context.Context, id string) option.Option[Session] {
	session, ok := s.store.get(id).Unwrap()
	if !ok {
		return option.None[Session]()
	}

	if time.Until(session.Token.Expires) < expirationLimit {
		token, err := s.discord.RefreshToken(ctx, session.Token.Refresh)
		if err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("failed to refresh token")
			s.store.delete(id)
			return option.None[Session]()
		}

		session.Token = token
		s.store.store(session)
	}

	return option.Some(session)
}

func (s *SessionStore) DeleteSession(ctx context.Context, id string) {
	s.store.delete(id)
}

type Session struct {
	ID    string
	Token Token
}

func (s *Session) Cookie() *http.Cookie {
	return &http.Cookie{
		Name:     "session",
		Value:    s.ID,
		Expires:  s.Token.Expires,
		HttpOnly: true,
	}
}
