package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/hadithopen-io/back/internal/auth/types"
	"github.com/hadithopen-io/back/pkg/pgscan"
	"github.com/hadithopen-io/back/pkg/usercontext"
)

type User struct{ db *sqlx.DB }

func NewUser(db *sqlx.DB) *User { return &User{db: db} }

type userWithPwd struct {
	ID      int64  `db:"id"`
	Email   string `db:"email"`
	Login   string `db:"login"`
	PwdHash string `db:"pwd_hash"`
}

func (u *User) GetByLogin(ctx context.Context, login string) (_ types.UserWithPwdHash, err error) {
	const query = `
select 
	id,
	email,
	login,
	pwd_hash
from public.users
where login = :login
`

	var user userWithPwd
	err = pgscan.Get(
		ctx,
		u.db,
		&user,
		query,
		map[string]any{
			"login": login,
		},
	)
	if err != nil {
		return types.UserWithPwdHash{}, err
	}

	return types.UserWithPwdHash{
		User: usercontext.User{
			ID:    user.ID,
			Email: user.Email,
			Login: user.Login,
		},
		PwdHash: user.PwdHash,
	}, nil
}
