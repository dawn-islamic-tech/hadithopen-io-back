package auth

import (
	"context"

	"github.com/hadithopen-io/back/internal/auth/types"
	"github.com/hadithopen-io/back/pkg/errors"
	"github.com/hadithopen-io/back/pkg/usercontext"
)

type Tokenizer interface {
	Tokenize(user usercontext.User) (
		token string,
		err error,
	)
}

type UserGetter interface {
	GetByLogin(ctx context.Context, login string) (
		user types.UserWithPwdHash,
		err error,
	)
}

type Encoder interface {
	Encode(from string) (
		hash string,
		err error,
	)
}

type Auth struct {
	user     UserGetter
	tokenize Tokenizer
	encoder  Encoder
}

func NewAuth(user UserGetter, tokenize Tokenizer, encoder Encoder) *Auth {
	return &Auth{
		user:     user,
		tokenize: tokenize,
		encoder:  encoder,
	}
}

func (a Auth) Login(ctx context.Context, pwd types.LoginWithPwd) (token string, err error) {
	user, err := a.user.GetByLogin(ctx, pwd.Login)
	if err != nil {
		return "", errors.Wrap(err, "getting user information")
	}

	hash, err := a.encoder.Encode(pwd.Pwd)
	if err != nil {
		return "", errors.Wrap(err, "decode prev info pwd")
	}

	if user.PwdHash != hash {
		return "", errors.New("you are input incorrect password")
	}

	return a.tokenize.Tokenize(
		user.User,
	)
}
