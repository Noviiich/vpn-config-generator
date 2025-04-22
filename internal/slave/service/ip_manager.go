package service

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/Noviiich/vpn-config-generator/internal/lib/e"
)

func (s *VPNService) getNextIP(ctx context.Context) (ip string, err error) {
	defer func() { err = e.WrapIfErr("can't get next ip: %v", err) }()
	lastIP, err := s.repo.GetIP(ctx)
	if err != nil {
		return "", err
	}

	ipFormat := net.ParseIP(lastIP)
	if ipFormat == nil {
		return "", err
	}

	ipParts := ipFormat.To4()
	if ipParts == nil {
		return "", err
	}

	lastOctet, _ := strconv.Atoi(fmt.Sprintf("%d", ipParts[3]))
	if lastOctet >= 254 { // Максимум x.x.x.254
		return "", err
	}

	newIP := fmt.Sprintf("%d.%d.%d.%d",
		ipParts[0],
		ipParts[1],
		ipParts[2],
		lastOctet+1,
	)

	if err = s.repo.UpdateIP(ctx, newIP); err != nil {
		return "", err
	}

	return newIP, nil
}
