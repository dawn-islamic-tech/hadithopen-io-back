package jwt

import (
	"context"

	"github.com/hadithopen-io/back/pkg/errors"
	"github.com/hadithopen-io/back/pkg/usercontext"
)

type Parser interface {
	Parse(token string) (
		user usercontext.User,
		err error,
	)
}

type Wrapper interface {
	WrapUser(parent context.Context, token string) (
		ctx context.Context,
		err error,
	)
}

type jwt struct {
	parser Parser
}

func NewWrapper(parser Parser) Wrapper {
	return &jwt{
		parser: parser,
	}
}

func (j *jwt) WrapUser(ctx context.Context, token string) (context.Context, error) {
	if token == "" {
		return ctx, errors.New("jwt: token has empty") // unauth
	}

	user, err := j.parser.Parse(token)
	if err != nil {
		return nil, errors.Wrap(err, "jwt: parse user from token")
	}

	return usercontext.Wrap(ctx, user), nil
}
