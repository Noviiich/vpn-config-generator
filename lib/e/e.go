package e

import (
	"errors"
	"fmt"
)

var (
	ErrUserNotFound        = errors.New("пользователя не существует")
	ErrUsersNotFound       = errors.New("пользователей не существует")
	ErrNotFound            = errors.New("неизвестная ошибка")
	ErrSubscriptionExpired = errors.New("ваша подписка истекла")
	ErrDevicesNotFound     = errors.New("у вас нет конфигураций")
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
