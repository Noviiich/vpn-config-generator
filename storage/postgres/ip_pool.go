package postgres

// func (s *Storage) GetIP(ctx context.Context, deviceID int) error {
// 	query := `SELECT INTO ip_pool (last_ip) VALUES ($1)`
// 	_, err := s.db.ExecContext(ctx, query, newIP)
// 	return err
// }

// func (s *Storage) UpdateIP(ctx context.Context, newIP string) error {
// 	query := `INSERT INTO ip_pool (last_ip) VALUES ($1)`
// 	_, err := s.db.ExecContext(ctx, query, newIP)
// 	return err
// }

// func (s *Storage) CreateIP(ctx context.Context) (string, error) {
// 	query := `SELECT last_ip FROM ip_pool ORDER BY id DESC LIMIT 1`

// 	var lastIP string
// 	err := s.db.QueryRowContext(ctx, query).Scan(&lastIP)

// 	if err != nil {
// 		// Если таблица пустая — возвращаем дефолтный IP
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return "10.0.0.1", nil
// 		}
// 		return "", err
// 	}

// 	return lastIP, nil
// }
