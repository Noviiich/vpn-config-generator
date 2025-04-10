package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/Noviiich/vpn-config-generator/internal/master/storage"
	"github.com/Noviiich/vpn-config-generator/lib/e"
)

func (s *VPNService) CreateAction(ctx context.Context, chatID int, typeID int) (err error) {
	user, err := s.GetUser(ctx, chatID)
	if err != nil {
		return err
	}

	action := &storage.Action{
		ActionID: typeID,
		UserID:   user.ID,
	}

	err = s.db.CreateAction(ctx, action)
	if err != nil {
		return e.ErrNotFound
	}

	return nil
}

func (s *VPNService) GetActions(ctx context.Context) (msg string, err error) {
	actions, err := s.db.GetActions(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", e.ErrActionsNotFound
		}
		return "", err
	}

	var acts []string
	for _, action := range actions {
		temp := string(action.ActionID) + "" + string(action.UserID)
		acts = append(acts, temp)
	}
	result := strings.Join(acts, "\n")
	return result, nil
}

func (s *VPNService) GetActionsByTime(ctx context.Context, since time.Time) ([]storage.Action, error) {
	actions, err := s.db.GetActionsByTime(ctx, since)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.ErrActionsNotFound
		}
		return nil, err
	}

	return actions, nil
}
