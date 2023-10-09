package errors

import (
	"errors"
	"fmt"
	"io"

	"go.uber.org/multierr"
)

type Error struct {
	msg string
}

func New(msg string) error { return &Error{msg: msg} }

func Is(err, target error) bool { return errors.Is(err, target) }

func (e *Error) Error() string { return e.msg }

func Unwrap(err error) error { return errors.Unwrap(err) }

func Errorf(format string, args ...any) error { return fmt.Errorf(format, args...) }

func Join(err ...error) error { return errors.Join(err...) }

func JoinCloser(err error, closer ...io.Closer) error {
	errs := make([]error, 0, len(closer)+1)
	errs = append(errs, err)
	for _, c := range closer {
		errs = append(errs, c.Close())
	}

	return Join(errs...)
}

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
