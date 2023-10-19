package middleware

import (
	"github.com/ogen-go/ogen/middleware"

	"github.com/hadithopen-io/back/pkg/empty"
	"github.com/hadithopen-io/back/pkg/errors"
	"github.com/hadithopen-io/back/pkg/usercontext"
)

func UserLang(req middleware.Request, next middleware.Next) (middleware.Response, error) {
	// According to past needs, there must already be a user
	u, err := usercontext.Get(req.Context)
	if err != nil && !errors.Is(err, usercontext.ErrCtxNotFound) {
		return middleware.Response{}, err
	}

	// If we have user in context, but lang param is empty, need set to default lang, for correctly app work.
	if u.Lang == "" {
		u.Lang = usercontext.DefaultLang
	}

	// For not logged-in user need made default.
	if !u.LoggedIn() {
		u = usercontext.Default()
	}

	req.Context = usercontext.Wrap(
		req.Context,
		usercontext.User{
			ID:    u.ID,
			Email: u.Email,
			Login: u.Login,

			// if lang == "" then use default user lang
			Lang: empty.Coalesce(
				req.Raw.URL.Query().Get("lang"),
				u.Lang,
			),
		},
	)

	return next(req)
}
