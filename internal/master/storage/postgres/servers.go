package postgres

import (
	"context"
	"database/sql"

	"github.com/Noviiich/vpn-config-generator/internal/master/storage"
)

func (s *Storage) GetServers(ctx context.Context) ([]storage.Server, error) {
	query := `SELECT id, name, country_flag, ip_address, port, protocol FROM servers`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servers []storage.Server
	for rows.Next() {
		var server storage.Server
		if err := rows.Scan(&server.ID, &server.Name, &server.CountryFlag, &server.IPAddress, &server.Port, &server.Protocol); err != nil {
			return nil, err
		}
		servers = append(servers, server)
	}

	return servers, nil
}

func (s *Storage) GetServer(ctx context.Context, serverID int) (*storage.Server, error) {
	query := `SELECT id, name, country_flag, ip_address, port, protocol FROM servers WHERE id = $1`
	row := s.db.QueryRowContext(ctx, query, serverID)

	var server storage.Server
	if err := row.Scan(&server.ID, &server.Name, &server.CountryFlag, &server.IPAddress, &server.Port, &server.Protocol); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Сервер с таким ID не найден
		}
		return nil, err
	}

	return &server, nil
}

func (s *Storage) AddServer(ctx context.Context, server *storage.Server) error {
	query := `INSERT INTO servers (name, country_flag, ip_address, port, protocol) VALUES ($1, $2, $3, $4, $5)`
	_, err := s.db.ExecContext(ctx, query, server.Name, server.CountryFlag, server.IPAddress, server.Port, server.Protocol)
	return err
}

func (s *Storage) DeleteServer(ctx context.Context, serverID int) error {
	query := `DELETE FROM servers WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, serverID)
	return err
}
