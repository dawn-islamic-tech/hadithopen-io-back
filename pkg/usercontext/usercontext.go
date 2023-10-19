package usercontext

import (
	"context"
	"fmt"

	"github.com/hadithopen-io/back/pkg/errors"
)

type User struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Login string `json:"login"`
	Lang  string `json:"lang"` // ru/en/ar/...

	authenticated bool
}

func Authenticated(u User) User {
	return User{
		ID:            u.ID,
		Email:         u.Email,
		Login:         u.Login,
		Lang:          u.Lang,
		authenticated: true,
	}
}

func (u User) LoggedIn() bool { return u.authenticated }

func (u User) String() string {
	return fmt.Sprintf("id - %d, email - %s, login - %s", u.ID, u.Email, u.Login)
}

type userContextKey struct{}

var ctxKey userContextKey

var ErrCtxNotFound = errors.New("from context user not found")

func Get(ctx context.Context) (User, error) {
	v := ctx.Value(ctxKey)
	if v == nil {
		return User{}, ErrCtxNotFound
	}

	u, _ := v.(User)
	return u, nil
}

func Wrap(ctx context.Context, user User) context.Context {
	return context.WithValue(
		ctx,
		ctxKey,
		user,
	)
}

const DefaultLang = "ru"

func Default() User {
	return User{
		Lang: DefaultLang,
	}
}
