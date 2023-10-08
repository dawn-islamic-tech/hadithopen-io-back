package middleware

import (
	"net/http"

	"github.com/ogen-go/ogen/middleware"

	"github.com/hadithopen-io/back/pkg/errors"
	"github.com/hadithopen-io/back/pkg/usercontext"
)

const cookieAuthName = "user"

func CookieAuth(req middleware.Request, next middleware.Next) (middleware.Response, error) {
	cook, err := getCookie(req.Raw)
	if err != nil {
		return middleware.Response{}, err
	}

	var user usercontext.User
	if cook == "" {
		// TODO:
		//  After developing the authorization module, you need to enrich the check and add a user from the cookie
		user = usercontext.Default()
	}

	req.Context = usercontext.Wrap(
		req.Context,
		user,
	)

	return next(req)
}

func getCookie(req *http.Request) (cookie string, _ error) {
	cook, err := req.Cookie(cookieAuthName)
	if errors.Is(err, http.ErrNoCookie) {
		return "", nil
	}

	return cook.Value, errors.Wrap(err, "getting cookie from http request")
}
