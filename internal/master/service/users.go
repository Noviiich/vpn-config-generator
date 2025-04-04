package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/Noviiich/vpn-config-generator/internal/master/storage"
	"github.com/Noviiich/vpn-config-generator/lib/e"
)

func (s *VPNService) GetUser(ctx context.Context, chatID int) (user *storage.User, err error) {
	user, err = s.db.GetUser(ctx, chatID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.ErrUserNotFound
		}
		return nil, e.ErrNotFound
	}

	return user, nil
}

func (s *VPNService) CreateUser(ctx context.Context, username string, chatID int) (err error) {
	user := &storage.User{
		TelegramID: chatID,
		Username:   username,
	}

	if err := s.db.CreateUser(ctx, user); err != nil {
		return e.ErrNotFound
	}

	return nil
}

func (s *VPNService) DeleteUser(ctx context.Context, chatID int) error {
	user, err := s.GetUser(ctx, chatID)
	if err != nil {
		return err
	}
	err = s.db.DeleteUser(ctx, user.ID)
	if err != nil {
		return e.ErrNotFound
	}
	return nil
}

func (s *VPNService) GetUsers(ctx context.Context) (string, error) {
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", e.ErrUserNotFound
		}
		return "", e.ErrNotFound
	}

	var usernames []string
	for _, user := range users {
		usernames = append(usernames, "@"+user.Username)
	}
	result := strings.Join(usernames, "\n")
	return result, nil
}
