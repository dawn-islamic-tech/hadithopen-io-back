package hs256

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/hadithopen-io/back/pkg/errors"
	"github.com/hadithopen-io/back/pkg/usercontext"
)

type HS256 struct {
	expiresAdd, refreshAdd time.Duration
	secret                 []byte
}

func NewHS256(expiresAdd, refreshAdd time.Duration, secret []byte) *HS256 {
	return &HS256{
		expiresAdd: expiresAdd,
		refreshAdd: refreshAdd,
		secret:     secret,
	}
}

type ClaimsWithRegister struct {
	usercontext.User
	jwt.RegisteredClaims
}

func (h *HS256) Tokenize(user usercontext.User) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		ClaimsWithRegister{
			User: user,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer: user.Login,
				ExpiresAt: jwt.NewNumericDate(
					time.Now().Add(h.expiresAdd),
				),
				IssuedAt: jwt.NewNumericDate(
					time.Now(),
				),
				NotBefore: jwt.NewNumericDate(
					time.Now(),
				),
			},
		},
	)

	signed, err := token.SignedString(h.secret)
	return signed, errors.Wrap(
		err,
		"make es256 signed token value",
	)
}

func (h *HS256) Parse(token string) (usercontext.User, error) {
	parsed, err := jwt.Parse(
		token,
		func(jt *jwt.Token) (any, error) {
			if _, ok := jt.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("undefined singing method")
			}

			return h.secret, nil
		},
	)
	if err != nil {
		return usercontext.User{}, errors.Wrap(err, "undefined")
	}
	if !parsed.Valid {
		return usercontext.User{}, errors.New("parsed token is not valid")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return usercontext.User{}, errors.New("undefined parsed claims")
	}

	id := int64From(claims["id"])
	if id == -1 {
		return usercontext.User{}, errors.New("undefined user id from parsed claims")
	}

	email, _ := claims["email"].(string)
	login, _ := claims["login"].(string)
	lang, _ := claims["lang"].(string)

	// If token parses OK, that use is authenticated in application.
	return usercontext.Authenticated(
		usercontext.User{
			ID:    id,
			Email: email,
			Login: login,
			Lang:  lang,
		},
	), nil
}

func int64From(v any) int64 {
	switch vv := v.(type) {
	case json.Number:
		i, _ := vv.Int64()
		return i

	case int:
		return int64(vv)

	case int32:
		return int64(vv)

	case float32:
		return int64(vv)

	case float64:
		return int64(vv)

	case int64:
		return vv

	case string:
		i, _ := strconv.ParseInt(vv, 10, 64)
		return i

	default:
		return -1
	}
}
