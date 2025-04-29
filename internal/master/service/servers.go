package service

import (
	"context"

	"github.com/Noviiich/vpn-config-generator/internal/master/storage"
)

func (s *VPNService) GetServers(ctx context.Context) ([]storage.Server, error) {
	return s.repo.GetServers(ctx)
}

func (s *VPNService) AddServer(ctx context.Context, server *storage.Server) error {
	return s.repo.AddServer(ctx, server)
}

func (s *VPNService) DeleteServer(ctx context.Context, serverID int) error {
	return s.repo.DeleteServer(ctx, serverID)
}
func (s *VPNService) GetServer(ctx context.Context, serverID int) (*storage.Server, error) {
	return s.repo.GetServer(ctx, serverID)
}
