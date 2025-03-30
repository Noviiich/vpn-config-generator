package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Noviiich/vpn-config-generator/lib/e"
	"github.com/Noviiich/vpn-config-generator/storage"
)

const (
	msgSubscriptionNotFound = `У вас нет подписки. 
Не расстраивайтесь, выполните команду /subscribe`
	msgSubscriptionActive = `Ваша подписка активна!!!
Она истекает через %d дней, %d часов`
)

// func (s *VPNService) UpdateSubscription(ctx context.Context, chatID int) (err error) {
// 	defer func() { err = e.WrapIfErr("can't update subscription", err) }()

// 	user, err := s.repo.GetUser(ctx, chatID)
// 	if err != nil {
// 		return err
// 	}

// 	user.SubscriptionActive = true
// 	user.SubscriptionExpiry = time.Now().Add(30 * 24 * time.Hour)

// 	return s.repo.UpdateUser(ctx, user)
// }

func (s *VPNService) CheckSubscription(ctx context.Context, chatID int) (msg string, err error) {
	user, err := s.GetUser(ctx, chatID)
	if err != nil {
		return "", err
	}

	sub, err := s.repo.CheckSubscription(ctx, user.ID)
	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return msgSubscriptionNotFound, nil
		}
		return "", e.ErrNotFound
	}

	return checkActive(sub)
}

func (s *VPNService) ActivateSubscription(ctx context.Context, chatID int) (err error) {
	user, err := s.GetUser(ctx, chatID)
	if err != nil {
		return err
	}

	err = s.repo.ActivateSubscription(ctx, user.ID, 1)
	if err != nil {
		log.Println(err)
		return e.ErrNotFound
	}

	return nil
}

func (s *VPNService) DeactivateSubscription(ctx context.Context, chatID int) (err error) {
	user, err := s.GetUser(ctx, chatID)
	if err != nil {
		return err
	}

	err = s.repo.DeactivateSubscription(ctx, user.ID)
	if err != nil {
		log.Println(err)
		return e.ErrNotFound
	}

	return nil
}

func checkActive(sub *storage.Subscription) (string, error) {
	if sub.IsActive {
		remaining := time.Until(sub.ExpiryDate)
		days := int(remaining.Hours()) / 24
		hours := int(remaining.Hours()) % 24

		message := fmt.Sprintf(msgSubscriptionActive, days, hours)
		return message, nil
	}

	return "", e.ErrSubscriptionExpired
}
