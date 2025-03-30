package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/Noviiich/vpn-config-generator/lib/e"
	"github.com/Noviiich/vpn-config-generator/storage"
)

func (s *VPNService) GetUser(ctx context.Context, chatID int) (user *storage.User, err error) {
	user, err = s.repo.GetUser(ctx, chatID)
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

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return e.ErrNotFound
	}

	return nil
}

func (s *VPNService) DeleteUser(ctx context.Context, chatID int) error {
	user, err := s.GetUser(ctx, chatID)
	if err != nil {
		return err
	}
	err = s.repo.DeleteUser(ctx, user.ID)
	if err != nil {
		return e.ErrNotFound
	}
	return nil
}

func (s *VPNService) isExistsUser(ctx context.Context, chatID int) (ex bool, err error) {
	exists, err := s.repo.IsExistsUser(ctx, chatID)
	if err != nil {
		return exists, e.Wrap("can't check if user exists", err)
	}
	return exists, nil
}

func (s *VPNService) GetUsers(ctx context.Context) (string, error) {
	users, err := s.repo.GetUsers(ctx)
	if err != nil {
		return "", e.Wrap("service: can't get users", err)
	}
	if users == nil {
		return "", e.ErrUsersNotFound
	}

	var usernames []string
	for _, user := range users {
		usernames = append(usernames, "@"+user.Username)
	}
	result := strings.Join(usernames, "\n")
	return result, nil
}
