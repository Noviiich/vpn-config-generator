package ports

import (
	"context"

	"github.com/Noviiich/vpn-config-generator/internal/core/domain"
)

type UserPort interface {
	All(context.Context) ([]domain.User, error)
	Create(context.Context, *domain.User) error
	Delete(context.Context, string)
	Get(context.Context, string) (domain.User, error)
}
