package usercontext

import (
	"context"

	"github.com/hadithopen-io/back/pkg/errors"
)

type User struct {
	Lang string `json:"lang"` // ru/en/ar/...
}

type userContextKey struct{}

var ctxKey userContextKey

func Get(ctx context.Context) (User, error) {
	v := ctx.Value(ctxKey)
	if v == nil {
		return User{}, errors.New("from context user not found")
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

func Default() User {
	const ruLang = "ru"
	return User{
		Lang: ruLang,
	}
}
