package sha256

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/hadithopen-io/back/pkg/errors"
)

type Encoder struct{ with string }

func NewEncoder(with string) *Encoder { return &Encoder{with: with} }

func (e Encoder) Encode(v string) (hash string, err error) {
	sh := sha256.New()

	if _, err := sh.Write([]byte(v)); err != nil {
		return "", errors.Wrap(err, "write hash")
	}

	if _, err := sh.Write([]byte(e.with)); err != nil {
		return "", errors.Wrap(err, "write with")
	}

	return hex.EncodeToString(
		sh.Sum(nil),
	), nil
}
