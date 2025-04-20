package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Noviiich/vpn-config-generator/lib/e"
)

func (s *VPNService) StatusSubscribtion(ctx context.Context, username string, chatID int) (st string, err error) {
	defer func() { err = e.WrapIfErr("can't get status subscription", err) }()

	exists, err := s.isExistsUser(ctx, chatID)
	if err != nil {
		return "", err
	}

	if !exists {
		err = s.CreateUser(ctx, username, chatID)
		if err != nil {
			return "", err
		}
	}

	user, err := s.repo.GetUser(ctx, chatID)
	if err != nil {
		return "", err
	}

	if user.SubscriptionActive {
		remaining := time.Until(user.SubscriptionExpiry)
		days := int(remaining.Hours()) / 24
		hours := int(remaining.Hours()) % 24

		msg := fmt.Sprintf(`Ваша подписка активна!!!
Она истекает через %d дней, %d часов`, days, hours)
		return msg, nil
	}

	return `У вас не подписки. Не расстраивайтесь, вы все еще можете ее оформить.
Для этого выполните /subscribe`, nil
}

func (s *VPNService) RemoveSubscription(ctx context.Context, chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't remove subscription", err) }()

	user, err := s.repo.GetUser(ctx, chatID)
	if err != nil {
		return err
	}

	user.SubscriptionActive = false

	return s.repo.UpdateUser(ctx, user)
}

func (s *VPNService) UpdateSubscription(ctx context.Context, chatID int, expiry time.Duration) (err error) {
	defer func() { err = e.WrapIfErr("can't update subscription", err) }()

	user, err := s.repo.GetUser(ctx, chatID)
	if err != nil {
		return err
	}

	user.SubscriptionActive = true
	user.SubscriptionExpiry = time.Now().Add(expiry)

	return s.repo.UpdateUser(ctx, user)
}
