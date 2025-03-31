package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Noviiich/vpn-config-generator/lib/e"
	"github.com/Noviiich/vpn-config-generator/storage"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func (s *VPNService) СreateDevice(ctx context.Context, userID int, subnet string) (err error) {
	// получение ip из вышедших, если они есть
	ipAddress, err := s.repo.GetIpIsNull(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	var ipUser string
	//если нет пустых ip из пула
	if errors.Is(err, sql.ErrNoRows) {
		//получение нового незадействованного ip
		ipUser, err = s.getNextIP(ctx, subnet)
		if err != nil {
			return err
		}
	}  else {
		ipUser = ipAddress.IP
		s.repo.UpdateIpPool(ctx, ipAddress.ID)
	}

	private, public, err := generateKeys()
	if err != nil {
		return err
	}

	device := &storage.Device{
		UserID:     userID,
		TypeId:     1,
		PrivateKey: private,
		PublicKey:  public,
		LastActive: time.Now(),
		IsActive:   true,
	}
	if err := s.repo.CreateDevice(ctx, device); err != nil {
		return err
	}



	return nil
}

func (s *VPNService) GetDevices(ctx context.Context, chatID int) (str string, err error) {
	devices, err := s.repo.GetDevices(ctx, chatID)
	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return "", e.ErrDevicesNotFound
		}
		return "", err
	}

	var usernames []string
	for _, device := range devices {
		usernames = append(usernames, strconv.Itoa(device.TypeId))
	}
	result := strings.Join(usernames, "\n")
	return result, nil
}

// func (s *VPNService) isExistsDevice(ctx context.Context, chatID int) (bool, error) {
// 	exists, err := s.repo.IsExistsDevice(ctx, chatID)
// 	if err != nil {
// 		return exists, e.Wrap("can't check if device exists", err)
// 	}
// 	return exists, nil
// }

func generateKeys() (string, string, error) {
	key, err := wgtypes.GenerateKey()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate private key: %w", err)
	}

	private := key.String()
	public := key.PublicKey().String()

	return private, public, nil
}
