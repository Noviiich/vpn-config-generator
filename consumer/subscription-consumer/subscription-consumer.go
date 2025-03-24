package subscription_consumer

import (
	"context"
	"log"
	"time"

	"github.com/Noviiich/vpn-config-generator/storage/postgres"
)

type Consumer struct {
	repo     *postgres.Storage
	interval time.Duration
}

func New(repo *postgres.Storage, interval time.Duration) *Consumer {
	return &Consumer{
		repo:     repo,
		interval: interval,
	}
}

func (c Consumer) Start() error {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	if err := c.checkSubscriptions(context.Background()); err != nil {
		log.Printf("[ERR] subscription consumer: %s", err.Error())
	}

	for range ticker.C {
		if err := c.checkSubscriptions(context.Background()); err != nil {
			log.Printf("[ERR] subscription consumer: %s", err.Error())
		}
	}
	return nil
}

func (c Consumer) checkSubscriptions(ctx context.Context) error {
	users, err := c.repo.GetUsers(ctx)
	if err != nil {
		log.Printf("[ERR] get users consumer: %s", err.Error())
	}

	now := time.Now()
	for _, user := range users {
		if user.SubscriptionExpiry.Before(now) {
			if err := c.repo.DeactivateSubscription(ctx, user.TelegramID); err != nil {
				log.Printf("Ошибка деактивации подписки для пользователя %d: %s", user.TelegramID, err.Error())
			} else {
				log.Printf("Подписка пользователя %d деактивирована", user.TelegramID)
			}
		}
	}

	return nil
}
