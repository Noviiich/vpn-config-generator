package service

import (
	"context"

	"github.com/Noviiich/vpn-config-generator/lib/e"
	"github.com/Noviiich/vpn-config-generator/storage"
)

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

func (s *VPNService) isExistsUser(ctx context.Context, chatID int) (ex bool, err error) {
	defer func() { err = e.WrapIfErr("can't is exists user", err) }()
	exists, err := s.repo.IsExistsUser(ctx, chatID)
	if err != nil {
		return exists, err
	}
	return exists, nil
}
