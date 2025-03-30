package e

import (
	"errors"
	"fmt"
)

var (
	ErrUserNotFound  = errors.New("пользователя не существует")
	ErrUsersNotFound = errors.New("пользователей не существует")
)

func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}

func WrapIfErr(msg string, err error) error {
	if err == nil {
		return nil
	}

	return Wrap(msg, err)
}
