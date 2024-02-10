package user

import (
	"context"

	"github.com/tarro-dev/discord-oauth/pkg/option"
)

type userKeyType string

const userKey userKeyType = "user"

func WithContext(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func Ctx(ctx context.Context) option.Option[*User] {
	user, ok := ctx.Value(userKey).(*User)
	if !ok {
		return option.None[*User]()
	}

	return option.Some(user)
}
