package types

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/hadithopen-io/back/pkg/errors"
)

type Translate struct {
	ID        int64  `db:"id"`
	Lang      string `db:"lang"`
	Translate string `db:"translate"`
}

type Translates struct {
	Values []Translate
}

func (t *Translates) Scan(src any) error {
	var bb []byte
	switch v := src.(type) {
	case nil:
		return nil

	case string:
		bb = []byte(v)

	case []byte:
		bb = v
	}

	var mp map[string]int64
	if err := json.Unmarshal(bb, &mp); err != nil {
		return errors.Wrap(err, "unmarshalling translates")
	}

	for k, v := range mp {
		t.Values = append(
			t.Values,
			Translate{
				ID:   v,
				Lang: k,
			},
		)
	}

	return nil
}

func (t Translates) Value() (driver.Value, error) {
	mp := make(map[string]int64, len(t.Values))
	for _, translate := range t.Values {
		mp[translate.Lang] = translate.ID
	}

	bb, err := json.Marshal(mp)
	return bb, errors.Wrap(
		err,
		"marshaling translates",
	)
}
