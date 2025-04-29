package config

type Config struct {
	Database Database
}

type Database struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func Load() *Config {
	return &Config{
		Database: Database{
			Host:     "postgres",
			Port:     "5432",
			Username: "novich",
			Password: "novich",
			DBName:   "vpndb",
			SSLMode:  "disable",
		},
	}
}
