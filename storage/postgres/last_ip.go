package postgres

import (
	"context"

	"github.com/Noviiich/vpn-config-generator/storage"
)

func (s *Storage) GetLastIP(ctx context.Context, subnet string) (*storage.LastIP, error) {
	query := `
			SELECT * FROM last_ip
			WHERE subnet = $1`
	var last storage.LastIP
	err := s.db.GetContext(ctx, &last, query, subnet)
	if err != nil {
		return nil, err
	}

	return &last, nil
}

func (s *Storage) UpdateIP(ctx context.Context, newIP string, subnet string) error {
	query := `
			UPDATE last_ip 
			SET ip = $1
			WHERE subnet = $2`
	_, err := s.db.ExecContext(ctx, query, newIP, subnet)
	return err
}
