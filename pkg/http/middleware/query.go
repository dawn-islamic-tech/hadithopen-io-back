package middleware

import (
	"github.com/ogen-go/ogen/middleware"

	"github.com/hadithopen-io/back/pkg/usercontext"
)

func QueryLang(req middleware.Request, next middleware.Next) (middleware.Response, error) {
	lang := req.Raw.URL.Query().Get("lang")
	if lang == "" {
		return next(req)
	}

	// According to past needs, there must already be a user
	_, err := usercontext.Get(req.Context)
	if err != nil {
		return middleware.Response{}, err
	}

	req.Context = usercontext.Wrap(
		req.Context,
		usercontext.User{
			Lang: lang,
		},
	)

	return next(req)
}
