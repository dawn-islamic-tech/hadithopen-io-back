package errors

import (
	"errors"
	"fmt"
	"io"

	"go.uber.org/multierr"
)

type Error struct {
	err error
}

func (e *Error) Error() string { return e.err.Error() }

func Unwrap(err error) error { return errors.Unwrap(err) }

func Errorf(format string, args ...any) error { return fmt.Errorf(format, args...) }

func Join(err ...error) error { return errors.Join(err...) }

func AppendCloser(into *error, closer io.Closer) { multierr.AppendInto(into, closer.Close()) }

func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%s: %w", msg, err)
}

func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf(
		"%s: %w",
		fmt.Sprintf(
			format,
			args...,
		),
		err,
	)
}
