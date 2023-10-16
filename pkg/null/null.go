package null

import (
	"database/sql"
	"reflect"
	"time"

	"github.com/hadithopen-io/back/pkg/errors"
)

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

type StringByte interface {
	~string | ~byte
}

type Time interface {
	time.Time
}

type Order interface {
	Number | StringByte | Time
}

type Null[T Order] struct {
	Val   T
	Valid bool
}

func (n *Null[T]) Scan(v any) (err error) {
	vv, ok, err := fromValue(n.Val, v)
	if err != nil {
		return err
	}

	n.Val = origin[T](vv)
	n.Valid = ok

	return nil
}

func origin[T Order](v any) T {
	var tt T
	out, _ := reflect.ValueOf(v).Convert(reflect.TypeOf(tt)).Interface().(T)
	return out
}

func fromValue(from, to any) (any, bool, error) {
	switch from.(type) {
	case int, int8, int16, int32, int64, float32, float64:
		return fromNumber(from, to)

	case bool:
		return nullBool(to)

	case string:
		return nullString(to)

	case byte:
		return nullByte(to)

	case time.Time:
		return nullTime(to)
	}

	switch reflect.TypeOf(from).Kind() {
	case reflect.Int:
		return nullInt[int](to)

	case reflect.Int8:
		return nullInt[int8](to)

	case reflect.Int16:
		return nullInt[int16](to)

	case reflect.Int32:
		return nullInt[int32](to)

	case reflect.Int64:
		return nullInt[int64](to)

	case reflect.String:
		return nullString(to)

	case reflect.Float64:
		return nullFloat[float64](to)

	case reflect.Float32:
		return nullFloat[float32](to)

	case reflect.Bool:
		return nullBool(to)
	}

	return nil, false, nil
}

func fromNumber(from, to any) (value any, valid bool, err error) {
	switch from.(type) {
	case int:
		value, valid, err = nullInt[int](to)

	case int8:
		value, valid, err = nullInt[int8](to)

	case int16:
		value, valid, err = nullInt[int16](to)

	case int64:
		value, valid, err = nullInt[int64](to)

	case float32:
		value, valid, err = nullFloat[float64](to)

	case float64:
		value, valid, err = nullFloat[float32](to)

	default:
		return nil, false, errors.New("undefined cast number type")
	}

	return value, valid, err
}

func nullInt[F int | int8 | int16 | int32 | int64](to any) (i F, valid bool, err error) {
	nn := &sql.NullInt64{}
	if err := nn.Scan(to); err != nil {
		return 0, false, errors.Wrap(err, "null: sql int64 scan")
	}

	return F(nn.Int64), nn.Valid, err
}

func nullFloat[F float32 | float64](to any) (f F, valid bool, err error) {
	nn := &sql.NullFloat64{}
	if err := nn.Scan(to); err != nil {
		return 0, false, errors.Wrap(err, "null: sql float64 scan")
	}

	return F(nn.Float64), nn.Valid, err
}

func nullBool(to any) (bool, bool, error) {
	nn := &sql.NullBool{}
	if err := nn.Scan(to); err != nil {
		return false, false, errors.Wrap(err, "null: sql bool scan")
	}

	return nn.Bool, nn.Valid, nil
}

func nullString(to any) (string, bool, error) {
	nn := &sql.NullString{}
	if err := nn.Scan(to); err != nil {
		return "", false, errors.Wrap(err, "null: sql string scan")
	}

	return nn.String, nn.Valid, nil
}

func nullByte(to any) (byte, bool, error) {
	nn := &sql.NullByte{}
	if err := nn.Scan(to); err != nil {
		return 0, false, errors.Wrap(err, "null: sql byte scan")
	}

	return nn.Byte, nn.Valid, nil
}

func nullTime(to any) (time.Time, bool, error) {
	nn := &sql.NullTime{}
	if err := nn.Scan(to); err != nil {
		return time.Time{}, false, errors.Wrap(err, "null: sql time scan")
	}

	return nn.Time, nn.Valid, nil
}
