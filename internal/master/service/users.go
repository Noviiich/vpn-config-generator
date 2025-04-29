package service

import (
	"context"

	"github.com/Noviiich/vpn-config-generator/internal/master/storage"
	"github.com/Noviiich/vpn-config-generator/internal/master/lib/e"
)

func (s *VPNService) GetUser(ctx context.Context, chatID int, username string) (user *storage.User, err error) {
	defer func() { err = e.WrapIfErr("can't get user", err) }()

	exists, err := s.isExistsUser(ctx, chatID)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = s.CreateUser(ctx, username, chatID)
		if err != nil {
			return nil, err
		}
	}

	user, err = s.repo.GetUser(ctx, chatID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *VPNService) CreateUser(ctx context.Context, username string, chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't create new user", err) }()
	user := &storage.User{
		TelegramID: chatID,
		Username:   username,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *VPNService) DeleteUser(ctx context.Context, chatID int) error {
	err := s.repo.DeleteUser(ctx, chatID)
	if err != nil {
		return e.Wrap("can't delete user", err)
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

func (s *VPNService) GetUsers(ctx context.Context, chatID int) ([]storage.User, error) {
	users, err := s.repo.GetUsers(ctx)
	if err != nil {
		return nil, e.Wrap("can't get users", err)
	}

	return users, nil
}
