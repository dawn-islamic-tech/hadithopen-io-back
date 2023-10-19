package types

import "github.com/hadithopen-io/back/pkg/usercontext"

type LoginWithPwd struct {
	Login, Pwd string
}

type UserWithPwdHash struct {
	PwdHash string

	usercontext.User
}
