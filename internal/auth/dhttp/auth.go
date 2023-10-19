package dhttp

import (
	"context"
	"net/http"

	"github.com/ogen-go/ogen/middleware"

	"github.com/hadithopen-io/back/internal/auth/dhttp/authgen"
	"github.com/hadithopen-io/back/internal/auth/types"
	"github.com/hadithopen-io/back/pkg/errors"
	"github.com/hadithopen-io/back/pkg/jwt"
)

type Auth interface {
	Login(ctx context.Context, pwd types.LoginWithPwd) (
		token string,
		err error,
	)
}

type AuthHandler struct {
	auth        Auth
	authWrapper jwt.Wrapper

	authgen.UnimplementedHandler
}

func NewAuthHandler(auth Auth, authWrapper jwt.Wrapper) *AuthHandler {
	return &AuthHandler{
		auth:        auth,
		authWrapper: authWrapper,
	}
}

func (a *AuthHandler) HandleCookieAuth(ctx context.Context, _ string, token authgen.CookieAuth) (
	context.Context,
	error,
) {
	return a.authWrapper.WrapUser(ctx, token.APIKey)
}

func (a *AuthHandler) Handler(m ...middleware.Middleware) (http.Handler, error) {
	handler, err := authgen.NewServer(
		a, // server implementation
		a, // secure middleware implementation
		authgen.WithMiddleware(
			m..., // global middlewares
		),
		authgen.WithPathPrefix(
			path, // basic prefix
		),
	)

	return handler, errors.Wrap(
		err,
		"make http hadith handler",
	)
}

const path = "/api/auth"

func (a *AuthHandler) Path() string { return path }

const cookieKey = "jwt"

func (a *AuthHandler) Login(ctx context.Context, req *authgen.UserLoginRequest) (
	*authgen.LoginOK,
	error,
) {
	token, err := a.auth.Login(
		ctx,
		types.LoginWithPwd{
			Login: req.Login,
			Pwd:   req.Pwd,
		},
	)
	if err != nil {
		return nil, err
	}

	return &authgen.LoginOK{
		SetCookie: authgen.NewOptString(
			cookieKey + "=" + token, // Info jwt need to front
		),
	}, nil
}

func (*AuthHandler) Logout(context.Context) (*authgen.LogoutOK, error) {
	return &authgen.LogoutOK{
		SetCookie: authgen.NewOptString(
			cookieKey + "=" + "", // Remove cookie value
		),
	}, nil
}
