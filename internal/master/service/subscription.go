package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Noviiich/vpn-config-generator/internal/master/storage"
	"github.com/Noviiich/vpn-config-generator/lib/e"
)

const (
	msgSubscriptionNotFound = `У вас нет подписки. 
Не расстраивайтесь, выполните команду /subscribe`
	msgSubscriptionActive = `Ваша подписка активна!!!
Она истекает через %d дней, %d часов`
)

func (s *VPNService) CreateSubscription(ctx context.Context, chatID int) (err error) {
	user, err := s.GetUser(ctx, chatID)
	if err != nil {
		return err
	}

	err = s.db.CreateSubscription(ctx, user.ID, 1)
	if err != nil {
		return e.ErrNotFound
	}

	return nil
}

func (s *VPNService) GetSubscription(ctx context.Context, chatID int) (msg string, err error) {
	user, err := s.GetUser(ctx, chatID)
	if err != nil {
		return "", err
	}

	sub, err := s.db.GetSubscription(ctx, user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return msgSubscriptionNotFound, nil
		}
		return "", e.ErrNotFound
	}

	return infoSubscription(sub)
}

func (s *VPNService) DeleteSubscription(ctx context.Context, chatID int) error {
	user, err := s.GetUser(ctx, chatID)
	if err != nil {
		return err
	}
	err = s.db.DeleteSubscription(ctx, user.ID)
	if err != nil {
		return e.ErrNotFound
	}

	return nil
}

func infoSubscription(sub *storage.Subscription) (string, error) {
	remaining := time.Until(sub.ExpiryDate)
	days := int(remaining.Hours()) / 24
	hours := int(remaining.Hours()) % 24

	message := fmt.Sprintf(msgSubscriptionActive, days, hours)
	return message, nil
}
