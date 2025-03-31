package service

import (
	"context"
	"fmt"
	"net"
	"strconv"
)

func (s *VPNService) getNextIP(ctx context.Context, subnet string) (string, error) {
	lastIP, err := s.repo.GetLastIP(ctx, "10.0.0.0/24")
	if err != nil {
		return "", err
	}

	ipFormat := net.ParseIP(lastIP.IP)
	if ipFormat == nil {
		return "", err
	}

	ipParts := ipFormat.To4()
	if ipParts == nil {
		return "", err
	}

	lastOctet, _ := strconv.Atoi(fmt.Sprintf("%d", ipParts[3]))
	if lastOctet >= 254 {
		return "", err
	}

	newIP := fmt.Sprintf("%d.%d.%d.%d",
		ipParts[0],
		ipParts[1],
		ipParts[3],
		lastOctet+1,
	)

	if err := s.repo.UpdateLastIP(ctx, newIP, subnet); err != nil {
		return "", err
	}

	return newIP, nil
}
