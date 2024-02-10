package auth

import (
	"context"

	"github.com/tarro-dev/discord-oauth/pkg/option"
)

type sessionKeyType string

const sessionKey sessionKeyType = "session"

func WithContext(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, sessionKey, session)
}

func Ctx(ctx context.Context) option.Option[*Session] {
	val := ctx.Value(sessionKey)
	if val == nil {
		return option.None[*Session]()
	}
	session, ok := ctx.Value(sessionKey).(*Session)
	if !ok {
		return option.None[*Session]()
	}

	return option.Some(session)
}
