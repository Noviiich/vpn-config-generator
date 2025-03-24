package subscription

import (
	"time"

	"github.com/Noviiich/vpn-config-generator/storage/postgres"
)

type Processor struct {
	repo     postgres.Storage
	interval time.Duration
}

func New(repo postgres.Storage) *Processor {
	return &Processor{
		repo: repo,
	}
}
